Review Dockerfile and the devcontainer.json files and update as needed.

NOTES:

1) note poor security posture of the dev container, which use
seccomp=unconfined and adds ptrace capability.
This is no doubt for debugging, "ease of use", but be aware of
the settings and update them if possible.
