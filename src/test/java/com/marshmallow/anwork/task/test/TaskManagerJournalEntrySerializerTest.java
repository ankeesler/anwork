package com.marshmallow.anwork.task.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.core.test.BaseSerializerTest;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskManagerActionType;
import com.marshmallow.anwork.task.TaskManagerJournalEntry;
import com.marshmallow.anwork.task.protobuf.TaskManagerActionTypeProtobuf;

import java.io.IOException;

import org.junit.Test;

public class TaskManagerJournalEntrySerializerTest
             extends BaseSerializerTest<TaskManagerJournalEntry> {

  private final TaskManager manager = new TaskManager();

  /**
   * Initialize a {@link TaskManagerJournalEntrySerializerTest} as a {@link BaseSerializerTest}.
   */
  public TaskManagerJournalEntrySerializerTest() {
    super(TaskManagerJournalEntry.SERIALIZER);
  }

  @Test
  public void singleTaskTest() throws IOException {
    manager.createTask("Task1", "This is task 1", 5);
    Task task = manager.getTasks()[0];
    TaskManagerJournalEntry entry
        = new TaskManagerJournalEntry(task, TaskManagerActionType.CREATE);
    TaskManagerJournalEntry deserializedEntry = runSerialization(entry);
    assertEquals(entry.getTitle(), deserializedEntry.getTitle());
    assertEquals(entry.getDescription(), deserializedEntry.getDescription());
    assertEquals(entry.getDate(), deserializedEntry.getDate());
  }

  @Test
  public void multiTaskTest() throws IOException {
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

    TaskManagerJournalEntry[] deserializedEntries = new TaskManagerJournalEntry[tasks.length];
    for (int i = 0; i < tasks.length; i++) {
      deserializedEntries[i] = runSerialization(entries[i]);
    }

    for (int i = 0; i < tasks.length; i++) {
      assertEquals(entries[i].getTitle(), deserializedEntries[i].getTitle());
      assertEquals(entries[i].getDescription(), deserializedEntries[i].getDescription());
      assertEquals(entries[i].getDate(), deserializedEntries[i].getDate());
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
  }
}
