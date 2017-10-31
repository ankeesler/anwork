package com.marshmallow.anwork.task;

import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.task.protobuf.TaskManagerJournalEntryProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskManagerJournalProtobuf;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * This is a {@link Journal} used by the {@link TaskManager}.
 *
 * <p>
 * Created Sep 29, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerJournal implements Journal<TaskManagerJournalEntry>,
                                           Serializable<TaskManagerJournalProtobuf> {

  public static final Serializer<TaskManagerJournal> SERIALIZER
      = new ProtobufSerializer<TaskManagerJournalProtobuf,
                               TaskManagerJournal>(() -> new TaskManagerJournal(),
                                                   TaskManagerJournalProtobuf.parser());

  private List<TaskManagerJournalEntry> entries = new ArrayList<TaskManagerJournalEntry>();

  @Override
  public void addEntry(TaskManagerJournalEntry entry) {
    entries.add(entry);
  }

  @Override
  public TaskManagerJournalEntry[] getEntries() {
    return entries.toArray(new TaskManagerJournalEntry[0]);
  }

  @Override
  public String toString() {
    return entries.toString();
  }

  /**
   * Get a {@link Journal} for this key. The key is the name of a task.
   *
   * <p>
   * If a {@link Task} has been deleted, then this method will return <code>null</code>.
   * </p>
   *
   * @param key The name of a task
   * @return A {@link Journal} for this key (i.e., task name)
   */
  public Journal<TaskManagerJournalEntry> filter(String key) {
    // TODO: cache me!
    TaskManagerJournal journal = new TaskManagerJournal();
    for (TaskManagerJournalEntry entry : entries) {
      if (entry.getTask().getName().equals(key)) {
        if (entry.getActionType().equals(TaskManagerActionType.DELETE)) {
          journal.entries.clear();
        } else {
          journal.entries.add(entry);
        }
      }
    }
    return (journal.entries.size() == 0 ? null : journal);
  }

  @Override
  public TaskManagerJournalProtobuf marshall() throws IOException {
    TaskManagerJournalProtobuf.Builder builder = TaskManagerJournalProtobuf.newBuilder();
    for (TaskManagerJournalEntry entry : entries) {
      builder.addEntries(entry.marshall());
    }
    return builder.build();
  }

  @Override
  public void unmarshall(TaskManagerJournalProtobuf t) throws IOException {
    for (int i = 0; i < t.getEntriesCount(); i++) {
      TaskManagerJournalEntryProtobuf entryProtobuf = t.getEntries(i);
      TaskManagerJournalEntry entry = new TaskManagerJournalEntry();
      entry.unmarshall(entryProtobuf);
      entries.add(entry);
    }
  }
}