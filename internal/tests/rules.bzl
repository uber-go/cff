load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("//rules:cff.bzl", "cff")
load("@bazel_skylib//lib:collections.bzl", "collections")

_default_deps = ["//src/go.uber.org/cff:go_default_library"]

_default_test_deps = [
    "@com_github_stretchr_testify//assert:go_default_library",
    "@com_github_stretchr_testify//require:go_default_library",
]

_base_importpath = "go.uber.org/cff/internal/tests/"

_visibility = ["//src/go.uber.org/cff:__subpackages__"]

def cff_internal_test(
        name,
        cff_srcs = None,
        test_srcs = None,
        srcs = None,
        deps = None,
        test_deps = None,
        instrument_all_tasks = False,
        online_scheduling = False):
    cff_srcs = cff_srcs or native.glob(include = ["*.go"], exclude = ["*_test.go"])
    test_srcs = test_srcs or native.glob(include = ["*_test.go"])

    srcs = srcs or []
    deps = deps or []
    test_deps = test_deps or []

    importpath = _base_importpath + name

    lib_name = name + "-library"
    test_name = name + "-test"
    cff_name = name + "-cff"

    deps += _default_deps
    test_deps += _default_test_deps

    deps = collections.uniq(deps)

    cff(
        name = cff_name,
        srcs = srcs,
        cff_srcs = cff_srcs,
        importpath = importpath,
        visibility = _visibility,
        deps = deps,
        instrument_all_tasks = instrument_all_tasks,
        online_scheduling = online_scheduling,
    )

    go_library(
        name = lib_name,
        srcs = [":" + cff_name],
        importpath = importpath,
        visibility = _visibility,
        deps = deps,
    )

    go_test(
        name = test_name,
        srcs = test_srcs + [":" + cff_name],
        race = "on",
        deps = collections.uniq(test_deps + deps),
    )
