import { TestBed, inject } from '@angular/core/testing';

import { HttpXsrfTokenExtractorDefault } from './http-xsrf-token-extractor.service';
import { SharedModule } from '../shared/shared.module';
import { CookieService } from "ngx-cookie";

describe('HttpXsrfTokenExtractorDefault', () => {
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
        HttpXsrfTokenExtractorDefault,
        { provide: CookieService, useValue: mockCookieService}
    ]
    });

  });

  it('should be initialized', inject([HttpXsrfTokenExtractorDefault], (service: HttpXsrfTokenExtractorDefault) => {
    expect(service).toBeTruthy();
  }));

  it('should be right get token', inject([HttpXsrfTokenExtractorDefault], (service: HttpXsrfTokenExtractorDefault) => {
    let token = service.getToken();
    expect(btoa(token)).toEqual("fdsa");
    mockCookieService.set(null);
    token = service.getToken();
    expect(token).toBeNull();
  }));

});
