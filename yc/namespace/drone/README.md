# Setup Drone CI and GitHub integration

According the [documentation](https://docs.drone.io/server/provider/github) you need to create a Drone CI application in GitHub and get client id and secret. Use https://drone.littlebit.space/login for redirect url.

After succesfull registration put the `client secret` in `yc/namespace/drone/server/github-secret.yaml` and `client id` into `DRONE_GITHUB_CLIENT_ID` in `yc/namespace/drone/server/service.yaml`.
