package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.event.EventLog;
import com.marshmallow.anwork.event.RamEventLog;
import com.marshmallow.anwork.task.LoggingTaskManager;
import com.marshmallow.anwork.task.TaskState;

import org.junit.Before;
import org.junit.Test;

/**
 * A test for the {@link LoggingTaskManager}.
 *
 * <p>
 * Created Sep 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class LoggingTaskManagerTest {

  private LoggingTaskManager manager = new LoggingTaskManager();
  private EventLog log = new RamEventLog();

  @Before
  public void setupManager() {
    manager.addLog(log);
  }

  @Test
  public void createTest() {
    manager.createTask("name-a", "description", 1);
    assertEquals(1, log.getEvents().length);
    manager.createTask("name-b", "description", 2);
    assertEquals(2, log.getEvents().length);
  }

  @Test
  public void deleteTest() {
    manager.createTask("name-a", "description", 1);
    manager.deleteTask("name-a");
    assertEquals(2, log.getEvents().length);
  }

  @Test
  public void setTaskTest() {
    manager.createTask("name-a", "description", 1);
    manager.setCurrentTask("name-a");
    assertEquals(2, log.getEvents().length);
    manager.setState("name-a", TaskState.WAITING);
    assertEquals(3, log.getEvents().length);
    manager.setState("name-a", TaskState.BLOCKED);
    assertEquals(4, log.getEvents().length);
  }

  @Test(expected = IllegalArgumentException.class)
  public void failureTest() {
    manager.createTask("name-a", "description", 1);
    assertEquals(1, log.getEvents().length);
    try {
      manager.createTask("name-a", "description", 2);
    } catch (IllegalArgumentException iae) {
      throw iae;
    } finally {
      assertEquals(2, log.getEvents().length);
    }
  }

  @Test
  public void removeLogTest() {
    assertTrue(manager.removeLog(log));
    manager.createTask("name", "description", 1);
    assertEquals(0, log.getEvents().length);
  }
}
