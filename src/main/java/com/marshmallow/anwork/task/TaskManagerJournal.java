package com.marshmallow.anwork.task;

import com.marshmallow.anwork.journal.Journal;

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
public class TaskManagerJournal implements Journal<TaskManagerJournalEntry> {

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
          journal.entries = new ArrayList<TaskManagerJournalEntry>();
        } else {
          journal.entries.add(entry);
        }
      }
    }
    return (journal.entries.size() == 0 ? null : journal);
  }
}