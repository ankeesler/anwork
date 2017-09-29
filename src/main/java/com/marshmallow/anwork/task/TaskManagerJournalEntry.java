package com.marshmallow.anwork.task;

import com.marshmallow.anwork.journal.JournalEntry;

import java.util.Date;

/**
 * This is a {@link JournalEntry} that records actions by the {@link TaskManager}.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerJournalEntry implements JournalEntry {

  private final Task task;
  private final TaskManagerActionType actionType;
  private final Date date;

  /**
   * Create a journal entry related to a {@link TaskManagerActionType} of action on a {@link Task}.
   *
   * @param task The {@link Task} on which the action of type {@link TaskManagerActionType} is
   *     taken
   * @param actionType The {@link TaskManagerActionType} of action that was taken on the
   *     {@link Task}
   */
  public TaskManagerJournalEntry(Task task, TaskManagerActionType actionType) {
    this.task = task;
    this.actionType = actionType;
    this.date = new Date();
  }

  @Override
  public String getTitle() {
    return String.format("%s:%s:%s", date, task.getName(), actionType.name());
  }

  @Override
  public String getDescription() {
    return getTitle();
  }

  @Override
  public Date getDate() {
    return date;
  }

  @Override
  public String toString() {
    return getTitle();
  }

  /**
   * Get the {@link Task} associated with this {@link TaskManagerJournalEntry}.
   *
   * @return The {@link Task} associated with this {@link TaskManagerJournalEntry}
   */
  public Task getTask() {
    return task;
  }

  /**
   * Get the {@link TaskManagerActionType} associated with this {@link TaskManagerJournalEntry}.
   *
   * @return The {@link TaskManagerActionType} associated with this
   *     {@link TaskManagerJournalEntry}.
   */
  public TaskManagerActionType getActionType() {
    return actionType;
  }
}
