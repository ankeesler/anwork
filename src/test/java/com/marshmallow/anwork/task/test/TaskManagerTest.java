package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;

import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

import org.junit.Test;

/**
 * This class tests the Task CRUD operations in a manager.
 *
 * <p>
 * Created Aug 29, 2017
 * </p>
 *
 * @author Andrew
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
    assertEquals(TaskState.WAITING, manager.getState("Task1"));

    manager.setState("Task1", TaskState.BLOCKED);
    assertEquals(TaskState.BLOCKED, manager.getState("Task1"));
    manager.setState("Task1", TaskState.BLOCKED);
    assertEquals(TaskState.BLOCKED, manager.getState("Task1"));
    manager.setState("Task1", TaskState.BLOCKED);
    assertEquals(TaskState.BLOCKED, manager.getState("Task1"));

    manager.createTask("Task2", "This is task 2.", 1);
    // See note about the waiting state above.
    assertEquals(TaskState.WAITING, manager.getState("Task2"));
    manager.setState("Task2", TaskState.RUNNING);
    assertEquals(TaskState.BLOCKED, manager.getState("Task1"));
    assertEquals(TaskState.RUNNING, manager.getState("Task2"));
  }

  @Test
  public void testTaskOrder() {
    manager.createTask("Task1", "This is task 1.", 2);
    Task[] tasks = manager.getTasks();
    assertEquals(1, tasks.length);
    assertEquals("Task1", tasks[0].getName());

    // Tasks are supposed to be returned in order of priority (see
    // TaskManager#getTasks).
    manager.createTask("Task2", "This is task 2.", 0);
    tasks = manager.getTasks();
    assertEquals(2, tasks.length);
    assertEquals("Task2", tasks[0].getName());
    assertEquals("Task1", tasks[1].getName());

    manager.createTask("Task3", "This is task 3.", 1);
    tasks = manager.getTasks();
    assertEquals(3, tasks.length);
    assertEquals("Task2", tasks[0].getName());
    assertEquals("Task3", tasks[1].getName());
    assertEquals("Task1", tasks[2].getName());

    manager.deleteTask("Task2");
    tasks = manager.getTasks();
    assertEquals(2, tasks.length);
    assertEquals("Task3", tasks[0].getName());
    assertEquals("Task1", tasks[1].getName());
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
    manager.setState("Task2", null);
  }
}
