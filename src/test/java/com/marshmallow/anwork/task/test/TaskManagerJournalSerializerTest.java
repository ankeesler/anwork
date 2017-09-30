package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.core.test.BaseSerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskManagerActionType;
import com.marshmallow.anwork.task.TaskManagerJournal;
import com.marshmallow.anwork.task.TaskManagerJournalEntry;

import java.io.IOException;

import org.junit.Test;

/**
 * A {@link BaseSerializerTest} for {@link TaskManagerJournal} objects.
 *
 * <p>
 * Created Sep 30, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerJournalSerializerTest extends BaseSerializerTest<TaskManagerJournal> {

  private final TaskManager manager = new TaskManager();
  private final TaskManagerJournal journal = new TaskManagerJournal();

  @Override
  public Serializer<TaskManagerJournal> getSerializer() {
    return TaskManagerJournal.SERIALIZER;
  }

  @Test
  public void testMultiEntryJournal() throws IOException {
    manager.createTask("Task0", "This is task 0", 0);
    manager.createTask("Task1", "This is task 1", 1);
    manager.createTask("Task2", "This is task 2", 2);
    Task[] tasks = manager.getTasks();
    assertEquals(3, tasks.length);

    TaskManagerJournalEntry[] entries = new TaskManagerJournalEntry[tasks.length];
    assertTrue(TaskManagerActionType.values().length >= tasks.length);
    for (int i = 0; i < tasks.length; i++) {
      entries[i]
          = new TaskManagerJournalEntry(tasks[i], TaskManagerActionType.values()[i]);
    }

    for (int i = 0; i < tasks.length; i++) {
      journal.addEntry(entries[i]);
    }

    TaskManagerJournal deserializedJournal = runSerialization(journal);
    TaskManagerJournalEntry[] deserializedEntries = deserializedJournal.getEntries();
    assertEquals(3, deserializedEntries.length);
    for (int i = 0; i < tasks.length; i++) {
      assertEquals(entries[i].getTitle(), deserializedEntries[i].getTitle());
      assertEquals(entries[i].getDescription(), deserializedEntries[i].getDescription());
      assertEquals(entries[i].getDate(), deserializedEntries[i].getDate());
    }
  }
}
