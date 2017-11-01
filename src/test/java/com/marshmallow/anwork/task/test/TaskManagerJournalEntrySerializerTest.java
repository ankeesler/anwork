package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.core.test.BaseSerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskManagerActionType;
import com.marshmallow.anwork.task.TaskManagerJournalEntry;
import com.marshmallow.anwork.task.protobuf.TaskManagerActionTypeProtobuf;

import java.io.IOException;

import org.junit.Test;

/**
 * A {@link BaseSerializerTest} for {@link TaskManagerJournalEntry} objects.
 *
 * <p>
 * Created Sep 30, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerJournalEntrySerializerTest
             extends BaseSerializerTest<TaskManagerJournalEntry> {

  private final TaskManager manager = new TaskManager();

  @Override
  protected Serializer<TaskManagerJournalEntry> getSerializer() {
    return TaskManagerJournalEntry.SERIALIZER;
  }

  @Test
  public void singleTaskTest() throws IOException {
    manager.createTask("Task1", "This is task 1", 5);
    Task task = manager.getTasks()[0];
    TaskManagerJournalEntry entry
        = new TaskManagerJournalEntry(task, TaskManagerActionType.CREATE, "hey");
    TaskManagerJournalEntry deserializedEntry = runSerialization(entry);
    assertEquals(entry.getTitle(), deserializedEntry.getTitle());
    assertEquals(entry.getDescription(), deserializedEntry.getDescription());
    assertEquals(entry.getDate(), deserializedEntry.getDate());
    assertEquals(entry.getDetail(), deserializedEntry.getDetail());
  }

  @Test
  public void multiTaskTest() throws IOException {
    manager.createTask("Task0", "This is task 0", 0);
    manager.createTask("Task1", "This is task 1", 1);
    manager.createTask("Task2", "This is task 2", 2);
    manager.createTask("Task3", "This is task 3", 2);
    Task[] tasks = manager.getTasks();
    assertEquals(4, tasks.length);

    TaskManagerJournalEntry[] entries = new TaskManagerJournalEntry[tasks.length];
    assertTrue(TaskManagerActionType.values().length >= tasks.length);
    for (int i = 0; i < tasks.length; i++) {
      entries[i] = new TaskManagerJournalEntry(tasks[i],
                                               TaskManagerActionType.values()[i],
                                               tasks[i].getName());
    }

    TaskManagerJournalEntry[] deserializedEntries = new TaskManagerJournalEntry[tasks.length];
    for (int i = 0; i < tasks.length; i++) {
      deserializedEntries[i] = runSerialization(entries[i]);
    }

    assertEquals(4, deserializedEntries.length);
    for (int i = 0; i < tasks.length; i++) {
      assertEquals(entries[i].getTitle(), deserializedEntries[i].getTitle());
      assertEquals(entries[i].getDescription(), deserializedEntries[i].getDescription());
      assertEquals(entries[i].getDate(), deserializedEntries[i].getDate());
      assertEquals(entries[i].getDetail(), deserializedEntries[i].getDetail());
    }
  }

  @Test
  public void testThatTaskManagerActionTypeEnumsLineUp() {
    assertEquals(TaskManagerActionType.CREATE.ordinal(),
                 TaskManagerActionTypeProtobuf.CREATE.ordinal());
    assertEquals(TaskManagerActionType.DELETE.ordinal(),
                 TaskManagerActionTypeProtobuf.DELETE.ordinal());
    assertEquals(TaskManagerActionType.SET_STATE.ordinal(),
                 TaskManagerActionTypeProtobuf.SET_STATE.ordinal());
    assertEquals(TaskManagerActionType.NOTE.ordinal(),
                 TaskManagerActionTypeProtobuf.NOTE.ordinal());
    assertEquals(TaskManagerActionType.SET_PRIORITY.ordinal(),
                 TaskManagerActionTypeProtobuf.SET_PRIORITY.ordinal());
  }
}
