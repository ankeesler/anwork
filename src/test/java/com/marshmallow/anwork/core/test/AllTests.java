package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.app.test.AllAppTests;
import com.marshmallow.anwork.journal.test.AllJournalTests;
import com.marshmallow.anwork.task.test.AllTaskTests;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;
import org.junit.runners.Suite.SuiteClasses;

@RunWith(Suite.class)
@SuiteClasses({
    AllCoreTests.class,
    AllTaskTests.class,
    AllJournalTests.class,
    AllAppTests.class,
  })
/**
 * This class is a wrapper around all of the unit tests for the ANWORK project.
 *
 * <p>
 * Created Aug 31, 2017
 * </p>
 *
 * @author Andrew
 */
public class AllTests { }
