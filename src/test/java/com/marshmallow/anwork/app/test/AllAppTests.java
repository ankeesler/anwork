package com.marshmallow.anwork.app.test;

import com.marshmallow.anwork.app.cli.test.CliTest;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

@RunWith(Suite.class)
@SuiteClasses({
    CliTest.class,
    AppTest.class,
  })
/**
 * This is a bucket for all application tests.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class AllAppTests { }
