load("@io_bazel_rules_go//go:def.bzl", "go_binary")
load("@rules_antlr//antlr:antlr4.bzl", "antlr")

antlr(
    name = "hdl",
    srcs = ["Hdl.g4"],
    language = "Go",
    package = "main",
)

go_binary(
    name = "compiler",
    srcs = [
        "compiler.go",
        ":hdl",
    ],
    deps = ["@com_github_antlr_antlr4//runtime/Go/antlr:go_default_library"],
)
