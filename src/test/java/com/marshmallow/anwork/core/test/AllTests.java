package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.app.test.AllAppTests;
import com.marshmallow.anwork.event.test.AllEventTests;
import com.marshmallow.anwork.task.test.AllTaskTests;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

@RunWith(Suite.class)
@SuiteClasses({
    AllCoreTests.class,
    AllTaskTests.class,
    AllEventTests.class,
    AllAppTests.class,
  })
/**
 * This class is a wrapper around all of the unit tests for the ANWORK project.
 *
 * @author Andrew
 * Created Aug 31, 2017
 */
public class AllTests { }
