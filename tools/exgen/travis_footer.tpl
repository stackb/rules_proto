before_install:
  - |
    if [ "${TRAVIS_OS_NAME}" = "osx" ]; then
      OS=darwin
    else
      sysctl kernel.unprivileged_userns_clone=1
      OS=linux
    fi
    wget -O install.sh "https://github.com/bazelbuild/bazel/releases/download/${BAZEL}/bazel-${BAZEL}-installer-${OS}-x86_64.sh"
    chmod +x install.sh
    ./install.sh --user
    rm -f install.sh

script:
  - |
    if [ -z "${COMMAND}" ] 
    then
      cd "${LANG}/example/${RULE}" && bazel build //:*
    else
      "${COMMAND}"
    fi

notifications:
  email: false
