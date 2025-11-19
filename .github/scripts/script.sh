act -j build -W .github/workflows/ci.yml \
    --container-architecture linux/amd64 \
    --privileged \
    -v /Users/chessy/.colima/default/docker.sock:/var/run/docker.sock
