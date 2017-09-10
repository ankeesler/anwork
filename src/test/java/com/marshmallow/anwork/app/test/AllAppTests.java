package com.marshmallow.anwork.app.test;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

import com.marshmallow.anwork.app.cli.test.CliTest;

@RunWith(Suite.class)
@SuiteClasses({
    CliTest.class,
    AppTest.class,
  })
/**
 * This is a bucket for all application tests.
 *
 * @author Andrew
 * Created Sep 9, 2017
 */
public class AllAppTests { }
