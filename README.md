 openssl req -new -subj "/C=GB/CN=foo" \
                  -addext "subjectAltName = DNS:localhost" \
                  -addext "certificatePolicies = 1.2.3.4" \
                  -newkey rsa:2048 -keyout key.pem -out req.pem




openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out example.crt \
            -keyout example.key



openssl req -newkey rsa:4096 \
-subj "/C=US/CN=localhost" \
-x509 \
-addext "subjectAltName = DNS:localhost" \
-sha256 \
-days 3650 \
-nodes \
-out example.crt \
-keyout example.key