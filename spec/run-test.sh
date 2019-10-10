#!/usr/bin/env roundup

describe "envy test"

before() {
  rm -rf tmp/spec/*
}

it_is_listed_in_help() {
  ($ENVY help 2>&1) | grep -w run
}

it_runs_run_section_wo_args() {
  TOKEN=foo ENVY_SECRET_PATH=spec/fixtures/secret $ENVY run spec/fixtures/config > tmp/spec/stdout
  diff tmp/spec/stdout spec/run-test.sh.out1
}

it_runs_specified_command() {
  expected="DATABASE_URL=postgres://pg_user:pg_password/server:5432/database/schema"
  actual=$(ENVY_SECRET_PATH=spec/fixtures/secret $ENVY run spec/fixtures/config env | grep DATABASE_URL)
  test $expected == $actual
}

it_overrides_environment() {
  expected="DATABASE_URL=postgres://pg_user:pg_password/server:5432/database/schema"
  actual=$(DATABASE_URL=nope ENVY_SECRET_PATH=spec/fixtures/secret $ENVY run spec/fixtures/config env | grep DATABASE_URL)
  test $expected == $actual
}
