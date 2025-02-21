package rotatesecretkey

import (
	"context"
	"fmt"
	"github.com/goharbor/harbor/src/common"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils"
	"github.com/goharbor/harbor/src/lib/config"
	configModels "github.com/goharbor/harbor/src/lib/config/models"
	"github.com/goharbor/harbor/src/lib/log"
	configDAO "github.com/goharbor/harbor/src/pkg/config/db/dao"
	oidcDAO "github.com/goharbor/harbor/src/pkg/oidc/dao"
	registryDAO "github.com/goharbor/harbor/src/pkg/reg/dao"
	"strings"
)

var (
	Ctl = NewController()
)

type Controller interface {
	// Rotate will load the sensitive data, tries to decrypt it using the current key and then re-encrypt it using the new key.
	// Then persist the new values to database, and set the system to readonly mode.
	Rotate(ctx context.Context, currentKey, newKey string, opt *Option) error
}

type controller struct{}

func (c controller) Rotate(ctx context.Context, currentKey, newKey string, opt *Option) error {
	log.Infof("Start rotating secrets")
	// Handle the configurations
	log.Infof("Re-encrypting configurations with the new key")
	var (
		updatedConfigKeys    []string
		updatedConfigEntries []configModels.ConfigEntry
	)

	configDAO := configDAO.New()
	configEntries, err := configDAO.GetConfigEntries(ctx)
	if err != nil {
		log.Errorf("Failed to get config entries: %v", err)
		return err
	}
	for _, e := range configEntries {
		if strings.HasPrefix(e.Value, utils.EncryptHeaderV1) {
			newValue, updated, err := reEncrypt(e.Value, currentKey, newKey)
			if err != nil {
				log.Errorf("Failed to re-encrypt configuration %s: %v", e.Key, err)
				return fmt.Errorf("failed to re-encrypt configuration %s: %v", e.Key, err)
			}
			if updated {
				updatedConfigKeys = append(updatedConfigKeys, e.Key)
				e.Value = newValue
				updatedConfigEntries = append(updatedConfigEntries,
					configModels.ConfigEntry{
						ID:    e.ID,
						Key:   e.Key,
						Value: newValue,
					})
			} else {
				log.Infof("Configuration %s is already encrypted with the new key", e.Key)
			}
		}
	}
	if len(updatedConfigEntries) > 0 {
		if err = configDAO.SaveConfigEntries(ctx, updatedConfigEntries); err != nil {
			log.Errorf("Failed to save config entries: %v", err)
			return err
		}
		log.Infof("Re-encrypted configurations : %v", updatedConfigKeys)
	}

	// Handle registries
	log.Infof("Re-encrypting registries' access secrets.")
	registryDAO := registryDAO.NewDAO()
	registries, err := registryDAO.List(ctx, nil)
	if err != nil {
		log.Errorf("Failed to get registries: %v", err)
		return err
	}
	for _, registry := range registries {
		newAccessSecret, updated, err := reEncrypt(registry.AccessSecret, currentKey, newKey)
		if err != nil {
			log.Errorf("Failed to re-encrypt registry %s access secret: %v", registry.Name, err)
			return fmt.Errorf("failed to re-encrypt registry %s access secret: %v", registry.Name, err)
		}
		if updated {
			registry.AccessSecret = newAccessSecret
			if err := registryDAO.Update(ctx, registry); err != nil {
				log.Errorf("Failed to update registry %s: %v", registry.Name, err)
				return fmt.Errorf("failed to update registry %s: %v", registry.Name, err)
			}
			log.Infof("Re-encrypted registry %s access secret", registry.Name)
		} else {
			log.Infof("Registry %s access secret is already encrypted with the new key", registry.Name)
		}
	}
	// Handle OIDC secrets
	if opt != nil && opt.SkipOIDCSecret {
		log.Infof("Skipped OIDC secret re-encryption. ")
	} else {
		log.Infof("Re-encrypting OIDC secrets")
		oidcDAO := oidcDAO.NewMetaDao()
		oidcUsers, err := oidcDAO.List(ctx, nil)
		if err != nil {
			log.Errorf("Failed to get OIDC users: %v", err)
			return err
		}
		for _, ou := range oidcUsers {
			if ou.Secret != "" {
				newSecret, updated, err := reEncrypt(ou.Secret, currentKey, newKey)
				if err != nil {
					log.Errorf("Failed to re-encrypt OIDC secret for user %d: %v", ou.UserID, err)
					return fmt.Errorf("failed to re-encrypt OIDC secret for user %d: %v", ou.UserID, err)
				}
				if !updated {
					log.Infof("OIDC secret for user %d is already encrypted with the new key", ou.UserID)
					continue
				}
				newToken, _, err := reEncrypt(ou.Token, currentKey, newKey)
				if err != nil {
					log.Errorf("Failed to re-encrypt token for user %d: %v", ou.UserID, err)
					return fmt.Errorf("failed to re-encrypt token for user %d: %v", ou.UserID, err)
				}
				if err := oidcDAO.Update(ctx, &models.OIDCUser{ID: ou.ID, Secret: newSecret, Token: newToken}, "secret", "token"); err != nil {
					log.Errorf("Failed to update OIDC secret for user %d: %v", ou.UserID, err)
					return fmt.Errorf("failed to update OIDC secret for user %d: %v", ou.UserID, err)
				}
				log.Infof("Re-encrypted OIDC secret for user %d", ou.UserID)
			}
		}
	}
	cfgMgr := config.GetCfgManager(ctx)
	cfgMgr.Set(ctx, common.ReadOnly, true)
	if err := cfgMgr.Save(ctx); err != nil {
		log.Errorf("Failed to set Harbor to ReadOnly after re-encrypting sensitive data: %v", err)
		return fmt.Errorf("failted to set Harbor to ReadOnly mode:%v", err)
	}
	log.Infof("Harbor is set to ReadOnly mode after re-encryption")
	return nil
}

func NewController() Controller {
	return &controller{}
}

type Option struct {
	// SkipOIDCSecret indicates whether to skip re-encrypting the OIDC secret
	SkipOIDCSecret bool
}

// tries to decrypt the input using the old key, and then re-encrypt it using the new key
// if the decryption fails, it will try to decrypt using the new key and return a flag indicating that the input was already encrypted using the new key
func reEncrypt(input, oldKey, newKey string) (string, bool, error) {
	plain, err := utils.ReversibleDecrypt(input, oldKey)
	if err != nil {
		// Check if it's already encrypted with the new key, if it is, return the input as is
		if _, err2 := utils.ReversibleDecrypt(input, newKey); err2 == nil {
			return input, false, nil
		} else {
			return "", false, err
		}
	}
	encrypted, err := utils.ReversibleEncrypt(plain, newKey)
	if err != nil {
		return "", false, err
	}
	return encrypted, true, nil
}
