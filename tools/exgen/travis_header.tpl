dist: trusty
sudo: required
osx_image: xcode8
language: android
android:
  components:
    - android-28
    - build-tools-28.0.3

os:
  - linux
  # - osx DISABLE FOR NOW UNTIL LINUX WORKING

env:
  - BAZEL=0.24.1 BAZEL_TEST='//example/routeguide:*'