package com.marshmallow.anwork.core.test;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

import com.marshmallow.anwork.event.test.EventLogTest;
import com.marshmallow.anwork.task.test.LoggingTaskManagerTest;
import com.marshmallow.anwork.task.test.TaskManagerTest;

@RunWith(Suite.class)
@SuiteClasses({
    SerializerTest.class,
    TaskManagerTest.class,
    EventLogTest.class,
    LoggingTaskManagerTest.class,
  })
/**
 * This class is a wrapper around all of the unit tests for the ANWORK project.
 *
 * @author Andrew
 * @date Aug 31, 2017
 */
public class AllTests { }
