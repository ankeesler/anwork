package com.marshmallow.anwork.task.test;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

@RunWith(Suite.class)
@SuiteClasses({
    LoggingTaskManagerTest.class,
    TaskManagerSerializerTest.class,
    TaskManagerTest.class,
    TaskSerializerTest.class,
  })
/**
 * This class is a wrapper around all of the task-related tests.
 *
 * @author Andrew
 * Created Sep 4, 2017
 */
public class AllTaskTests { }
