import { TestBed, inject } from '@angular/core/testing';

import { HttpXsrfTokenExtractorToBeUsed } from './http-xsrf-token-extractor.service';
import { SharedModule } from '../shared/shared.module';
import { CookieService } from "ngx-cookie";

describe('HttpXsrfTokenExtractorToBeUsed', () => {
    let cookie =  "fdsa|ds";
  let mockCookieService =  {
      get: function () {
          return cookie;
      },
      set: function (cookieStr: string) {
          cookie = cookieStr;
      }
  };
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule
      ],
      providers: [
        HttpXsrfTokenExtractorToBeUsed,
        { provide: CookieService, useValue: mockCookieService}
    ]
    });

  });

  it('should be initialized', inject([HttpXsrfTokenExtractorToBeUsed], (service: HttpXsrfTokenExtractorToBeUsed) => {
    expect(service).toBeTruthy();
  }));

  it('should be right get token', inject([HttpXsrfTokenExtractorToBeUsed], (service: HttpXsrfTokenExtractorToBeUsed) => {
    let token = service.getToken();
    expect(btoa(token)).toEqual("fdsa");
    mockCookieService.set(null);
    token = service.getToken();
    expect(token).toBeNull();
  }));

});
