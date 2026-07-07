# syntax=docker/dockerfile:1.10

# Static-file host for the ESDM schemas. Serves the embedded YAML
# files at https://schema.esdm.io/<name>/v<N> with the right MIME
# type, conservative caching, and permissive CORS so editors and
# language servers can fetch them from anywhere. TLS termination is
# expected at the ingress layer; this image listens on plain HTTP/80.

FROM nginx:1.29.3-alpine3.22

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY core /usr/share/nginx/html/core
COPY domain-storytelling /usr/share/nginx/html/domain-storytelling
COPY given-when-then /usr/share/nginx/html/given-when-then
