chromaprint-server
===

A server where you can send audio files, and you get back chromaprints.

Useful if you have a  _lot_ of audio to fingerprint and you want to split up the tasks over many servers

# Dependencies

ffmpeg and libchromaprint

On debain you would run:

`apt-get install ffmpeg libchromaprint0 libchromaprint-dev` ( -dev is only needed if you are building )

# Example usage:

## PNG output

```
ben@metropolis:/tmp$ curl -sv --data-binary @test.mp3 -H 'Content-Type: audio/mp3' 'localhost:6464/chromaprint?png=y' -o png.png
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to localhost (127.0.0.1) port 6464 (#0)
> POST /chromaprint?png=y HTTP/1.1
> User-Agent: curl/7.35.0
> Host: localhost:6464
> Accept: */*
> Content-Type: audio/mp3
> Content-Length: 1544192
> Expect: 100-continue
> 
< HTTP/1.1 100 Continue
} [data not shown]
< HTTP/1.1 200 OK
< Content-Type: image/png
< Date: Wed, 04 Jan 2017 21:30:04 GMT
< Content-Length: 1751
< 
{ [data not shown]
* Connection #0 to host localhost left intact
```

Returns a PNG like so:

![exampleout](exampleoutput.png)


## Base64 output

```
ben@metropolis:/tmp$ curl --data-binary @test.mp3 -H 'Content-Type: audio/mp3' 'localhost:6464/chromaprint' 
AQAB8moUKYuSaMGlHGV-pOqPXoJ-XMd1HNof-Jj8I09ExErEoySxjzr848mJScdx5UQXwj4ehR98aOyPo46PrTt-iMdDB1aPOsel7Hh2NCe0Ey86
6lAl6UKZE33GHelR_3gyPNExZeTQs0BPcKmELVIW5Ghf6GMq_MhjPDVRO2iO39iRZwejG38DUXRE6FF4-EGPR0P4wx2PfMcOUWFhRtHRH3mUozzC
nMah3YU15UIfaDz-wye-wFWPF9bRQyfqHTLTY_fxw7KQ58end_B1hOonnCO0rCz6B_0VNCehHuVh0cN3fNByPBp6uNWF59CRKnjGo9mOXajUo1lu
jCei59DS8eh34fhROcGVHMelHd_RaDouHXdxKM-Fz0fDW6gkhKuOJg6PMjmRbsnA-Al41IedHJUSzUaXUwncd-hTXLvwHY8gjseFZ0QPkcdlPLi4
Dx6D_sIJTUWzHA93VDrE4zme5fCPL2h-HOV8GG0UGjWaF18SolFlhO_BTESYx1CoEkeYRiwmmTxyHY_whAt0fchzKcMfnCiyRLiPJw0uWsaD9nCH
E18OK7rRH5rR7NglDtoUPYIjoy_C54X7oQ-uoYVXDztOcJYQJ8qC_TC_QjtzPMiD5xnuHD3RPDd2IR-45dCWNwncS0OpBGcp_Gg0QbuOH_qCHiau
CA-ufbBS9LKg6fCGljxqfYcj4YcOK_OOw1VQajR-9Edzoo1OhHnxo04SV2iOvGAmHoqS38h7PEbTYb_xRwgv6FmWI4_y4E9wGl-QJUKhG1cS5sNV
6HiYo9bRiDmuoLvQHLVyDQ2V4sqH5ugDKSzcsEU3FX2G5sfF4lEeNNuhgz9i4qgSFj0eFRfT4DvOBndy4-FxHhuVbsT140C00oeeITy-z2h_-PvQ
I7wKMWmObx0RxiHOHDfCPJjOHLqOXMSU_XgOJ-JxvsIlE1cvpA-SHb_xHM0TVEcV7cMUl8SPHA2fQ3sk_HiO_hD94dSOw6Jk_GiPloFo8UB-uDwO
R-XxyEW5Hyf6wz7yQcu0XTgv_EFaZRleHPUUF030H_9RcTz8azhhkkevozeaE3mH5zKqajman3jx469wXngcEfYPnXNw9NGJTyIOvceD_PB09MJN
_NiPFxPvIjfEYXoiI9KPG9cSGT1THD18Hu2y4PIR_ngDMSfS8uiHc3iHqV3h63A0jXgF50f44SiP99iLlkefwXlw4tDh6TgOndhaoloH5MWP_nA6
FN_RwxwsX3jw44R_fMuNHOfxCxPxj8MHHf2L9PiDH6IAUIAQpyxihhgMBBJCWKAARggAA5gBEhgGDYOOQIIAQcoAIQAAikhCjBRESwSdMMYwIgRQ
WADHlGAGSILAQIAgaAARAhDRFDKCUAEIE0AI5AQCQgFnILICCGbYMU4gAQRyAlghCTLAQGEcAAwJ4oQQBmFABQCEGCSUMoIjZggSzjIDFVNCKGOE
IcYIYSwCTkEBnCHMOMIMklIYB4iiAAlhACDQCccEIBIhhJQAQhAhDEMAISBMUgIiYASwQgFDJGACGQSAUggIg4wBQgAEGFLMAQEIFxRBJIBCwDOB
BDAGIKM4NwpBAoxgwBFBQEGaIAAAEAgBYswRwBijAHMCAIWAsAgoA6xhSAUCiGdOCUIEAFQA4QEyQAknDFPEEOIUIwAIgAABChmglABEAAEJMIIg
xgBwBACGgAMWEIIQcEoAAZQRQlBAiGNKIEEUUUgYwCxlgghgiGNFDIUoIMpQJCgAVQkAKFFCGICYEtwZAQI 
```

## Raw input

If you already have data in the raw format ( signed 16 bit little ed ints, mono, 44100Hz ) you can skip the
ffmpeg decoder:

```
ben@metropolis:/tmp$ curl --data-binary @test.raw 'localhost:6464/chromaprint?raw=y'
```

