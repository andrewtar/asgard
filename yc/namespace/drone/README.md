# Setup Drone CI and Gitea integration

According the [documentation](https://docs.drone.io/server/provider/gitea/) you need to create a Drone CI application in Gitea and get client id and secret. [Documentation](https://docs.drone.io/server/provider/gitea/). Use https://drone.littlebit.space/login for redirect url.

After succesfull registration put the `client secret` in `yc/namespace/drone/server/gitea-secret.yaml` and `client id` into `DRONE_GITEA_CLIENT_ID` in `yc/namespace/drone/server/service.yaml`.
