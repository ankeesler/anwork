package com.marshmallow.anwork.core.test;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

import com.marshmallow.anwork.event.test.AllEventTests;
import com.marshmallow.anwork.task.test.AllTaskTests;

@RunWith(Suite.class)
@SuiteClasses({
    AllTaskTests.class,
    AllEventTests.class,
    FilePersisterTest.class,
  })
/**
 * This class is a wrapper around all of the unit tests for the ANWORK project.
 *
 * @author Andrew
 * @date Aug 31, 2017
 */
public class AllTests { }
