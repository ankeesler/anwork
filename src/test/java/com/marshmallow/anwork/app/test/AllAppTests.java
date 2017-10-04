package com.marshmallow.anwork.app.test;

import com.marshmallow.anwork.app.cli.test.BasicCliTest;
import com.marshmallow.anwork.app.cli.test.CliXmlTest;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

@RunWith(Suite.class)
@SuiteClasses({
    BasicCliTest.class,
    CliXmlTest.class,
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
