package com.marshmallow.anwork.task.test;

import org.junit.Test;

import com.marshmallow.anwork.task.TaskManager;

import static org.junit.Assert.*;

public class TaskManagerTest {

  private static final String CONTEXT = "tuna";

  @Test
  public void testCreateTask() {
    TaskManager manager = new TaskManager(CONTEXT);
    assertEquals(0, manager.getTaskCount());

    manager.createTask("Task1", "This is task 1.", 1);
    assertEquals(1, manager.getTaskCount());

    manager.createTask("Task2", "This is task 2.", 1);
    assertEquals(2, manager.getTaskCount());
  }

  @Test
  public void testDeleteTask() {
    TaskManager manager = new TaskManager(CONTEXT);
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
    TaskManager manager = new TaskManager(CONTEXT);
    assertNull(manager.getCurrentTask());
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setCurrentTask("Task1");
    assertEquals("Task1", manager.getCurrentTask());
  }

  @Test
  public void testState() {
    TaskManager manager = new TaskManager(CONTEXT);
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
    TaskManager manager = new TaskManager(CONTEXT);
    manager.createTask("Task1", "This is task 1.", 1);
    manager.createTask("Task1", "This is another task 1.", 1);
  }

  @Test(expected = IllegalArgumentException.class)
  public void deleteUnknownTask() {
    TaskManager manager = new TaskManager(CONTEXT);
    manager.createTask("Task1", "This is task 1.", 1);
    manager.deleteTask("Task2");
  }

  @Test(expected = IllegalArgumentException.class)
  public void setUnknownCurrentTask() {
    TaskManager manager = new TaskManager(CONTEXT);
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setCurrentTask("Task2");
  }

  @Test(expected = IllegalArgumentException.class)
  public void setBadState() {
    TaskManager manager = new TaskManager(CONTEXT);
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setState("Task2", "whatever");
  }
}
