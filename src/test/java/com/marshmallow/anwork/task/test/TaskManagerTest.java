package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotEquals;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.journal.JournalEntry;
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
    assertEquals(0, manager.getTasks().length);
    assertEquals(0, manager.getJournal().getEntries().length);

    manager.createTask("Task1", "This is task 1.", 1);
    assertEquals(1, manager.getTasks().length);
    assertJournalEntriesEqual("Task1");

    manager.createTask("Task2", "This is task 2.", 1);
    assertEquals(2, manager.getTasks().length);
    assertNotEquals(manager.getTasks()[0].getId(), manager.getTasks()[1].getId());
    assertJournalEntriesEqual("Task1", "Task2");

    assertJournalEntrySize("Task1", 1);
    assertJournalEntrySize("Task2", 1);
    assertJournalEntrySize("this task does not exist", null);
  }

  @Test
  public void testDeleteTask() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.createTask("Task2", "This is task 2.", 1);

    manager.deleteTask("Task1");
    assertEquals(1, manager.getTasks().length);

    manager.createTask("Task3", "This is task 3.", 1);
    assertEquals(2, manager.getTasks().length);
    assertJournalEntriesEqual("Task1", "Task2", "Task1", "Task3");
    assertJournalEntrySize("Task1", null);
    assertJournalEntrySize("Task2", 1);
    assertJournalEntrySize("Task3", 1);

    manager.deleteTask("Task2");
    assertEquals(1, manager.getTasks().length);

    manager.deleteTask("Task3");
    assertEquals(0, manager.getTasks().length);

    assertJournalEntriesEqual("Task1", "Task2", "Task1", "Task3", "Task2", "Task3");
    assertJournalEntrySize("Task1", null);
    assertJournalEntrySize("Task2", null);
    assertJournalEntrySize("Task3", null);
  }

  @Test
  public void testState() {
    manager.createTask("Task1", "This is task 1.", 1);

    // By default, tasks start out in the WAITING state.
    assertEquals(TaskState.WAITING, manager.getState("Task1"));

    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(1, entries.length);
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

    assertJournalEntriesEqual("Task1", "Task1", "Task1", "Task1", "Task2", "Task2");
    assertJournalEntrySize("Task1", 4);
    assertJournalEntrySize("Task2", 2);
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

  @Test
  public void testAddAndRemoveTaskSameName() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setState("Task1", TaskState.RUNNING);
    manager.setState("Task1", TaskState.FINISHED);
    assertJournalEntriesEqual("Task1", "Task1", "Task1");
    assertJournalEntrySize("Task1", 3);

    manager.deleteTask("Task1");
    assertJournalEntriesEqual("Task1", "Task1", "Task1", "Task1");
    assertJournalEntrySize("Task1", null);

    manager.createTask("Task1", "This is task 1, again.", 2);
    assertJournalEntriesEqual("Task1", "Task1", "Task1", "Task1", "Task1");
    assertJournalEntrySize("Task1", 1);
  }

  @Test
  public void testAddSimilarlyNamedTasks() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.createTask("ask", "This is ask.", 2);
    manager.setState("ask", TaskState.RUNNING);
    manager.setState("Task1", TaskState.BLOCKED);
    manager.setState("ask", TaskState.BLOCKED);
    assertJournalEntriesEqual("Task1", "ask", "ask", "Task1", "ask");
    assertJournalEntrySize("Task1", 2);
    assertJournalEntrySize("ask", 3);
  }

  @Test
  public void testAddingNotes() {
    manager.createTask("Task1", "", 1);
    manager.createTask("Task2", "", 2);

    manager.addNote("Task1", "Note1");
    assertJournalEntriesEqual("Task1", "Task2", "Task1");
    assertJournalEntrySize("Task1", 2);
    assertJournalEntrySize("Task2", 1);

    manager.addNote("Task2", "Note2");
    assertJournalEntriesEqual("Task1", "Task2", "Task1", "Task2");
    assertJournalEntrySize("Task1", 2);
    assertJournalEntrySize("Task2", 2);

    manager.addNote("Task1", "Note3");
    assertJournalEntriesEqual("Task1", "Task2", "Task1", "Task2", "Task1");
    assertJournalEntrySize("Task1", 3);
    assertJournalEntrySize("Task2", 2);
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
  public void setBadState() {
    manager.createTask("Task1", "This is task 1.", 1);
    manager.setState("Task2", null);
  }

  @Test(expected = IllegalArgumentException.class)
  public void addBadNote() {
    manager.createTask("Task1", "", 1);
    manager.addNote("Task2", "hey");
  }

  // This method asserts that the provided list of expectedNames matches what is actually in the
  // manager's journal.
  private void assertJournalEntriesEqual(String...expectedNames) {
    // See TaskManagerJournalCache for more discussion on this weird logic.
    JournalEntry[] entries = manager.getJournal().getEntries();
    assertEquals(expectedNames.length, entries.length);
    for (int i = 0; i < expectedNames.length; i++) {
      assertTrue(entries[i].getTitle().contains(expectedNames[i]));
    }
  }

  // This method asserts that the journal for the provided key is of length size.
  // You can pass null for size to indicate that there is no journal for this key.
  private void assertJournalEntrySize(String key, Integer size) {
    Journal<?> journal = manager.getJournal(key);
    if (size == null) {
      assertNull(journal);
    } else {
      assertEquals(size.intValue(), journal.getEntries().length);
    }
  }
}
