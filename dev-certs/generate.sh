#!/bin/bash

rm *.csr *.key *.pem *.srl *.conf

echo "Generating CA key + cert"
openssl req -x509 -newkey rsa:4096 -nodes -keyout cakey.key -subj "/CN=TestCA/C=MY" -days 3650 -out cacert.pem

echo "Generating server cert"
cat > csr.conf <<EOF
[ req ]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = MY
CN = localhost

[ req_ext ]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
IP.1 = 127.0.0.1

EOF

openssl req -new -sha256 -nodes -out csr.csr -newkey rsa:4096 -keyout server.key  -config csr.conf

openssl x509 -req -in csr.csr -CA cacert.pem -CAkey cakey.key -CAcreateserial -out server.pem -days 9000 -extfile csr.conf -extensions req_ext


echo "Generating client cert"
cat > csrclient.conf <<EOF
[ req ]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = MY
CN = client

[ req_ext ]
keyUsage = keyEncipherment
extendedKeyUsage = clientAuth

EOF


openssl req -new -sha256 -nodes -out client.csr -newkey rsa:4096 -keyout client.key  -config csrclient.conf

openssl x509 -req -in client.csr -CA cacert.pem -CAkey cakey.key -CAcreateserial -out client.pem -days 9000 -extfile csrclient.conf -extensions req_ext


echo "Generating test cert"
cat > csrtest.conf <<EOF
[ req ]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C = MY
CN = test

[ req_ext ]
keyUsage = keyEncipherment
extendedKeyUsage = clientAuth

EOF

openssl req -new -sha256 -nodes -out test.csr -newkey rsa:4096 -keyout test.key  -config csrtest.conf

openssl x509 -req -in test.csr -CA cacert.pem -CAkey cakey.key -CAcreateserial -out test.pem -days 9000 -extfile csrtest.conf -extensions req_ext


echo 'generating unknown certificate'
openssl req -x509 -newkey rsa:4096 -nodes -keyout unknown.key -out unknown.pem -days 3650 \
		-subj "/C=XX/ST=unknown/L=unknown/O=unknown/CN=unknown"
