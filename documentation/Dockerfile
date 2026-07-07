# syntax=docker/dockerfile:1.10

# The .esdm.yaml documents under documentation/guides/<slug>/model/ are
# validated by `esdm lint` in the QA workflow before merge to main. This image
# build trusts main to be green and does not re-run the lint.

FROM squidfunk/mkdocs-material:9.7.6 AS build-documentation

RUN pip install "mkdocs-material[imaging]"

ADD . .

RUN mkdocs build

# ---------------------------------------------------------

FROM nginx:1.29.3-alpine3.22

COPY --from=build-documentation /docs/site /usr/share/nginx/html
