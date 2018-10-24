# Hiboot - web/cli application framework 

<p align="center">
  <a href="https://hiboot.hidevops.io">
    <img src="https://github.com/hidevopsio/hiboot/blob/master/hiboot.png?raw=true" alt="hiboot">
  </a>
</p>

<p align="center">
  <a href="https://travis-ci.org/hidevopsio/hiboot?branch=master">
    <img src="https://travis-ci.org/hidevopsio/hiboot.svg?branch=master" alt="Build Status"/>
  </a>
  <a href="https://codecov.io/gh/hidevopsio/hiboot">
    <img src="https://codecov.io/gh/hidevopsio/hiboot/branch/master/graph/badge.svg" />
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
      <img src="https://img.shields.io/badge/License-Apache%202.0-green.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/hidevopsio/hiboot">
      <img src="https://goreportcard.com/badge/github.com/hidevopsio/hiboot" />
  </a>
  <a href="https://godoc.org/github.com/hidevopsio/hiboot">
      <img src="https://godoc.org/github.com/golang/gddo?status.svg" />
  </a>
  <a href="https://gitter.im/hidevopsio/hiboot">
      <img src="https://img.shields.io/badge/GITTER-join%20chat-green.svg" />
  </a>
</p>

## About

Hiboot is a cloud native web and cli application framework written in Go.

Hiboot is not trying to reinvent everything, it integrates the popular libraries but make them simpler, easier to use. It borrowed some of the Spring features like dependency injection, aspect oriented programming, and auto configuration. You can integrate any other libraries easily by auto configuration with dependency injection support.

If you are a Java developer, you can start coding in Go without learning curve.

## Overview

* Web MVC (Model-View-Controller).
* Auto Configuration, pre-create instance with properties configs for dependency injection.
* Dependency injection with struct tag name **\`inject:""\`**, **Constructor** func, or **Method**.

## Getting Started

* [Hiboot Documentations](https://hiboot.hidevops.io) https://hiboot.hidevops.io
* [Hiboot 文档](https://hiboot.hidevops.io/cn) https://hiboot.hidevops.io/cn

## Community Contributions Guide

Thank you for considering contributing to the Hiboot framework, The contribution guide can be found [here](CONTRIBUTING.md).

## License

© John Deng, 2017 ~ time.Now

Released under the [Apache License 2.0](https://github.com/hidevopsio/hiboot/blob/master/LICENSE)