package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;

import com.marshmallow.anwork.core.test.BaseSerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;

import java.io.IOException;

import org.junit.Test;

/**
 * A {@link BaseSerializerTest} for {@link TaskManager} objects.
 *
 * <p>
 * Created Sep 4, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerSerializerTest extends BaseSerializerTest<TaskManager> {

  private static class TaskInfo {
    private final String name;
    private final String description;
    private final int priority;

    public TaskInfo(String name, String description, int priority) {
      this.name = name;
      this.description = description;
      this.priority = priority;
    }

    public String getName() {
      return name;
    }

    public String getDescription() {
      return description;
    }

    public int getPriority() {
      return priority;
    }
  }

  private static void assertTaskManagerEquals(TaskManager manager, TaskInfo...infos) {
    Task[] tasks = manager.getTasks();
    assertEquals(infos.length, tasks.length);
    for (int i = 0; i < infos.length; i++) {
      assertEquals(infos[i].getName(), tasks[i].getName());
      assertEquals(infos[i].getDescription(), tasks[i].getDescription());
      assertEquals(infos[i].getPriority(), tasks[i].getPriority());
    }
  }

  private TaskManager manager = new TaskManager();

  /**
   * Instantiate this test as a {@link BaseSerializerTest}.
   */
  public TaskManagerSerializerTest() {
    super(TaskManager.SERIALIZER);
  }

  @Test
  public void testNoTasks() throws IOException {
    manager = runSerialization(manager);
    assertTaskManagerEquals(manager);
  }

  @Test
  public void testSingleTask() throws IOException {
    manager.createTask("task-a", "This is task a", 1);
    manager = runSerialization(manager);

    assertTaskManagerEquals(manager, new TaskInfo[] {
        new TaskInfo("task-a", "This is task a", 1),
    });
  }

  @Test
  public void testMultipleTasks() throws IOException {
    manager.createTask("task-a", "This is task a", 1);
    manager.createTask("task-b", "This is task b", 2);
    manager.createTask("task-c", "This is task c", 0);
    manager = runSerialization(manager);

    assertTaskManagerEquals(manager, new TaskInfo[] {
        new TaskInfo("task-c", "This is task c", 0),
        new TaskInfo("task-a", "This is task a", 1),
        new TaskInfo("task-b", "This is task b", 2),
    });
  }

  @Test
  public void testTasksWithSpaces() throws IOException {
    manager.createTask("task a", "this is task a", 1);
    manager.createTask("t a s k   b", "this is task b", 2);
    manager = runSerialization(manager);

    Task[] tasks = manager.getTasks();
    assertEquals(2, tasks.length);

    assertEquals("task a", tasks[0].getName());
    assertEquals("this is task a", tasks[0].getDescription());
    assertEquals(1, tasks[0].getPriority());

    assertEquals("t a s k   b", tasks[1].getName());
    assertEquals("this is task b", tasks[1].getDescription());
    assertEquals(2, tasks[1].getPriority());

    assertTaskManagerEquals(manager, new TaskInfo[] {
        new TaskInfo("task a", "this is task a", 1),
        new TaskInfo("t a s k   b", "this is task b", 2),
    });
  }

  @Test
  public void testALotOfTasks() throws IOException {
    for (int i = 0; i < 100; i++) {
      manager.createTask("task-" + i, "this is task " + i, i);
    }

    manager = runSerialization(manager);

    Task[] tasks = manager.getTasks();
    assertEquals(100, tasks.length);
    for (int i = 0; i < 100; i++) {
      assertEquals("task-" + i, tasks[i].getName());
      assertEquals("this is task " + i, tasks[i].getDescription());
      assertEquals(i, tasks[i].getPriority());
    }
  }

  @Test
  public void testLargeTasks() throws IOException {
    String message
        = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus sed metus non nunc"
          + " varius porttitor. Vivamus ut bibendum eros, nec rutrum nulla. Vestibulum suscipita"
          + " dolor quis auctor. Vestibulum sodales quis velit ut mollis. Phasellus mattis tempor"
          + " arcu et efficitur. Donec ac elit efficitur, facilisis ligula et, tristique sapien."
          + " Mauris feugiat ante sed accumsan sollicitudin. Aliquam laoreet eros urna, nec"
          + " efficitur metus aliquam facilisis. Ut interdum nec dolor id dictum. Nam fringilla"
          + " pulvinar ex placerat scelerisque. Nullam vestibulum mi eget risus malesuada volutpat."
          + " Mauris pulvinar risus et faucibus faucibus. Nulla facilisi. Duis non nunc nibh."
          + " Integer orci odio, blandit cursus ornare a, molestie ut massa. Integer ultricies"
          + " rutrum enim, euismod posuere mauris luctus vel.";
    manager.createTask("this is task 1", message, 1);
    manager.createTask("this is task 2", message, 2);
    manager.createTask("this is task 3", message, 3);

    manager = runSerialization(manager);

    assertTaskManagerEquals(manager, new TaskInfo[] {
        new TaskInfo("this is task 1", message, 1),
        new TaskInfo("this is task 2", message, 2),
        new TaskInfo("this is task 3", message, 3),
    });
  }
}
