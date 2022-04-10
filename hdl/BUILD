load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@rules_antlr//antlr:antlr4.bzl", "antlr")

antlr(
    name = "hdl",
    srcs = ["Hdl.g4"],
    language = "Go",
    package = "reader",
    visitor = True,
)

go_library(
    name = "reader",
    srcs = [
        "reader.go",
        ":hdl",
    ],
    importpath = "github.com/dubov94/es-computer/hdl/reader",
    deps = ["@com_github_antlr_antlr4//runtime/Go/antlr:go_default_library"],
)

go_binary(
    name = "compiler",
    srcs = ["compiler.go"],
    deps = [":reader"],
)
