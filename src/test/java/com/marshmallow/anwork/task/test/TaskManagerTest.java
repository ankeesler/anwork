package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;

import com.marshmallow.anwork.task.TaskManager;

import org.junit.Test;

/**
 * This class tests the Task CRUD operations in a manager.
 *
 * @author Andrew
 * Created Aug 29, 2017
 */
public class TaskManagerTest {

  private TaskManager manager = new TaskManager();

  @Test
  public void testCreateTask() {
    assertEquals(0, manager.getTaskCount());

    manager.createTask("Task1", "This is task 1.", 1);
    assertEquals(1, manager.getTaskCount());

    manager.createTask("Task2", "This is task 2.", 1);
    assertEquals(2, manager.getTaskCount());
  }

  @Test
  public void testDeleteTask() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.createTask("Task2", "This is task 2.", 1);

    manager.deleteTask("Task1");
    assertEquals(1, manager.getTaskCount());

    manager.createTask("Task3", "This is task 3.", 1);
    assertEquals(2, manager.getTaskCount());

    manager.deleteTask("Task2");
    assertEquals(1, manager.getTaskCount());

    manager.deleteTask("Task3");
    assertEquals(0, manager.getTaskCount());
  }

  @Test
  public void testCurrentTask() {
    assertNull(manager.getCurrentTask());
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setCurrentTask("Task1");
    assertEquals("Task1", manager.getCurrentTask());
  }

  @Test
  public void testState() {
    manager.createTask("Task1", "This is task 1.", 1);

    // By default, tasks start out in the WAITING state.
    assertEquals("waiting", manager.getState("Task1"));

    manager.setState("Task1", "blocked");
    assertEquals("blocked", manager.getState("Task1"));
    manager.setState("Task1", "BLOCKED");
    assertEquals("blocked", manager.getState("Task1"));
    manager.setState("Task1", "BlOcKeD");
    assertEquals("blocked", manager.getState("Task1"));

    manager.createTask("Task2", "This is task 2.", 1);
    // See note about the waiting state above.
    assertEquals("waiting", manager.getState("Task2"));
    manager.setState("Task2", "running");
    assertEquals("blocked", manager.getState("Task1"));
    assertEquals("running", manager.getState("Task2"));
  }

  @Test(expected = IllegalArgumentException.class)
  public void testDuplicateTask() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.createTask("Task1", "This is another task 1.", 1);
  }

  @Test(expected = IllegalArgumentException.class)
  public void deleteUnknownTask() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.deleteTask("Task2");
  }

  @Test(expected = IllegalArgumentException.class)
  public void setUnknownCurrentTask() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setCurrentTask("Task2");
  }

  @Test(expected = IllegalArgumentException.class)
  public void setBadState() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setState("Task2", "whatever");
  }
}
