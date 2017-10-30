package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotEquals;

import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.core.test.BaseSerializerTest;
import com.marshmallow.anwork.journal.JournalEntry;
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

  private static void assertTaskManagerEquals(TaskManager manager,
                                              Task[] expectedTasks,
                                              JournalEntry[] expectedEntries) {
    Task[] tasks = manager.getTasks();
    assertEquals(expectedTasks.length, tasks.length);
    for (int i = 0; i < expectedTasks.length; i++) {
      assertEquals(expectedTasks[i].getName(), tasks[i].getName());
      assertEquals(expectedTasks[i].getDescription(), tasks[i].getDescription());
      assertEquals(expectedTasks[i].getPriority(), tasks[i].getPriority());
    }

    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(expectedEntries.length, entries.length);
    for (int i = 0; i < expectedEntries.length; i++) {
      assertEquals(expectedEntries[i].getTitle(), entries[i].getTitle());
      assertEquals(expectedEntries[i].getDescription(), entries[i].getDescription());
      assertEquals(expectedEntries[i].getDate(), entries[i].getDate());
    }
  }

  private TaskManager manager = new TaskManager();

  @Override
  protected Serializer<TaskManager> getSerializer() {
    return TaskManager.SERIALIZER;
  }

  @Test
  public void testNoTasks() throws IOException {
    manager = runSerialization(manager);
    assertTaskManagerEquals(manager, new Task[0], new JournalEntry[0]);
  }

  @Test
  public void testSingleTask() throws IOException {
    manager.createTask("task-a", "This is task a", 1);

    Task[] tasks = manager.getTasks();
    assertEquals(1, tasks.length);
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(1, entries.length);

    manager = runSerialization(manager);
    assertTaskManagerEquals(manager, tasks, entries);
  }

  @Test
  public void testMultipleTasks() throws IOException {
    manager.createTask("task-a", "This is task a", 1);
    manager.createTask("task-b", "This is task b", 2);
    manager.createTask("task-c", "This is task c", 0);
    Task[] tasks = manager.getTasks();
    assertEquals(3, tasks.length);
    assertNotEquals(tasks[0].getId(), tasks[1].getId());
    assertNotEquals(tasks[0].getId(), tasks[2].getId());
    assertNotEquals(tasks[1].getId(), tasks[2].getId());
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(3, entries.length);

    manager = runSerialization(manager);
    assertTaskManagerEquals(manager, tasks, entries);
  }

  @Test
  public void testTasksWithSpaces() throws IOException {
    manager.createTask("task a", "this is task a", 1);
    manager.createTask("t a s k   b", "this is task b", 2);
    Task[] tasks = manager.getTasks();
    assertEquals(2, tasks.length);
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(2, entries.length);

    manager = runSerialization(manager);
    assertTaskManagerEquals(manager, tasks, entries);
  }

  @Test
  public void testALotOfTasks() throws IOException {
    for (int i = 0; i < 100; i++) {
      manager.createTask("task-" + i, "this is task " + i, i);
    }
    Task[] tasks = manager.getTasks();
    assertEquals(100, tasks.length);
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(100, entries.length);

    manager = runSerialization(manager);
    assertTaskManagerEquals(manager, tasks, entries);
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
    Task[] tasks = manager.getTasks();
    assertEquals(3, tasks.length);
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(3, entries.length);

    manager = runSerialization(manager);
    assertTaskManagerEquals(manager, tasks, entries);
  }
}
