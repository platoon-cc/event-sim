#!/bin/env bash

export VERSION=$(cat cmd/.version)
echo Adding git tag with version v${VERSION}
git tag v${VERSION}
git push origin v${VERSION}
