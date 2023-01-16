load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "099a9fb96a376ccbbb7d291ed4ecbdfd42f6bc822ab77ae6f1b5cb9e914e94fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "efbbba6ac1a4fd342d5122cbdfdb82aeb2cf2862e35022c752eaddffada7c3f3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.27.0/bazel-gazelle-v0.27.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_rules_dependencies()

go_register_toolchains(version = "1.18.3")

go_repository(
    name = "com_github_go_telegram_bot_api",
    importpath = "github.com/go-telegram-bot-api/telegram-bot-api/v5",
    sha256 = "80b3ded8f36a050aca8101322ddee721fefe285af59a8560a1c059b961526cc2",
    strip_prefix = "telegram-bot-api-5.5.1",
    type = "zip",
    urls = ["https://github.com/go-telegram-bot-api/telegram-bot-api/archive/refs/tags/v5.5.1.zip"],
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sha256 = "89fc6ed61fd9e01732a77fa4649b975ba113e8d111a23a6ec8457c5811407f48",
    strip_prefix = "testify-1.8.1",
    type = "zip",
    urls = ["https://github.com/stretchr/testify/archive/refs/tags/v1.8.1.zip"],
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sha256 = "d2db56f11ee5180d52b78b3f542fb538d7bfb8ce6d6a7b26d7722c38a3f0810a",
    strip_prefix = "go-spew-1.1.1",
    type = "zip",
    urls = ["https://github.com/davecgh/go-spew/archive/refs/tags/v1.1.1.zip"],
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sha256 = "d003bda4f967bc70af7ab710067763a58016dfb5a82aa00e451c5e457af99555",
    strip_prefix = "go-difflib-5d4384ee4fb2527b0a1256a821ebfc92f91efefc",
    type = "zip",
    urls = ["https://github.com/pmezard/go-difflib/archive/5d4384ee4fb2527b0a1256a821ebfc92f91efefc.zip"],
)

go_repository(
    name = "in_gopkg_yaml_v3",
    importpath = "gopkg.in/yaml.v3",
    sha256 = "a2cb35f011c3579f1b959a07f754f272d8bb373bd89d7f3f609cd3e50298046f",
    strip_prefix = "yaml-3.0.1",
    type = "zip",
    urls = ["https://github.com/go-yaml/yaml/archive/refs/tags/v3.0.1.zip"],
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sha256 = "c21f6fffef4f58f3d8b17d88053f326d58b6b979ce71ceca986a10691bfdd6d5",
    strip_prefix = "logrus-1.9.0",
    type = "zip",
    urls = ["https://github.com/sirupsen/logrus/archive/refs/tags/v1.9.0.zip"],
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sha256 = "06e72b7166ec73e4ee6454e25d68abba7bf1e0beabebe09fe38b379b836dc976",
    strip_prefix = "sys-0.3.0",
    type = "zip",
    urls = ["https://github.com/golang/sys/archive/refs/tags/v0.3.0.zip"],
)

go_repository(
    name = "com_github_golang_jwt_jwt",
    importpath = "github.com/golang-jwt/jwt/v4",
    sha256 = "572c3adef54ab0dc9e2cf8f1b92865ea1d44a4301db73953a5c3e95c42f05087",
    strip_prefix = "jwt-4.4.3",
    type = "zip",
    urls = ["https://github.com/golang-jwt/jwt/archive/refs/tags/v4.4.3.zip"],
)

# rules_docker
# 0.23 breaks something around cc toolchains
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "59536e6ae64359b716ba9c46c39183403b01eabfbd57578e84398b4829ca499a",
    strip_prefix = "rules_docker-0.22.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.22.0/rules_docker-v0.22.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//repositories:deps.bzl",
    container_deps = "deps",
)

container_deps()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)

# https://console.cloud.google.com/gcr/images/distroless/global/base
container_pull(
    name = "distroless-container-image",
    digest = "sha256:b9b124f955961599e72630654107a0cf04e08e6fa777fa250b8f840728abd770",
    registry = "gcr.io",
    repository = "distroless/base",
)

load(
    "@io_bazel_rules_docker//toolchains/docker:toolchain.bzl",
    docker_toolchain_configure = "toolchain_configure",
)

docker_toolchain_configure(name = "docker_config")

# Import toolchain repositories for remote executions, but register the
# toolchains using --extra_toolchains on the command line to get precedence.
local_repository(
    name = "remote_config_cc",
    path = "tools/remote-toolchains/ubuntu-act-22-04/local_config_cc",
)
