#!/bin/sh

openssl pkcs12 -legacy -in ftp_u1.p12 -clcerts -nokeys -out $KEY_NAME1.crt -passin pass:$KEY1_PASS
openssl pkcs12 -legacy -in ftp_u1.p12 -nocerts -out $KEY_NAME1.key -nodes -passin pass:$KEY1_PASS
openssl pkcs12 -legacy -in ftp_i2.p12 -clcerts -nokeys -out $KEY_NAME1.crt -passin pass:$KEY2_PASS
openssl pkcs12 -legacy -in ftp_i2.p12 -nocerts -out $KEY_NAME1.key -nodes -passin pass:$KEY2_PASS
