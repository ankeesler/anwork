package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;

import com.marshmallow.anwork.core.test.SerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

import java.io.IOException;

import org.junit.Test;

/**
 * A {@link SerializerTest} for {@link TaskManager} objects.
 *
 * <p>
 * Created Sep 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerSerializerTest extends SerializerTest<TaskManager> {

  private TaskManager manager = new TaskManager();

  /**
   * Instantiate this test as a {@link SerializerTest}.
   */
  public TaskManagerSerializerTest() {
    super(TaskManager.SERIALIZER);
  }

  @Test
  public void testNoTasks() throws IOException {
    manager = runSerialization(manager);
    assertEquals(0, manager.getTasks().length);
  }

  @Test
  public void testSingleTask() throws IOException {
    manager.createTask("task-a", "This is task a", 1);
    manager = runSerialization(manager);

    Task[] tasks = manager.getTasks();
    assertEquals(1, tasks.length);
    assertEquals("task-a", tasks[0].getName());
    assertEquals("This is task a", tasks[0].getDescription());
    assertEquals(1, tasks[0].getPriority());
  }

  @Test
  public void testMultipleTasks() throws IOException {
    manager.createTask("task-a", "This is task a", 1);
    manager.createTask("task-b", "This is task b", 2);
    manager.createTask("task-c", "This is task c", 0);
    manager = runSerialization(manager);

    Task[] tasks = manager.getTasks();
    assertEquals(3, tasks.length);

    assertEquals("task-c", tasks[0].getName());
    assertEquals("This is task c", tasks[0].getDescription());
    assertEquals(0, tasks[0].getPriority());

    assertEquals("task-a", tasks[1].getName());
    assertEquals("This is task a", tasks[1].getDescription());
    assertEquals(1, tasks[1].getPriority());

    assertEquals("task-b", tasks[2].getName());
    assertEquals("This is task b", tasks[2].getDescription());
    assertEquals(2, tasks[2].getPriority());
  }
}
