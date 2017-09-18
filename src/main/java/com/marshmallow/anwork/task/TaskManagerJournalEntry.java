package com.marshmallow.anwork.task;

import com.marshmallow.anwork.journal.BaseJournalEntry;
import com.marshmallow.anwork.journal.JournalEntry;

/**
 * This is a {@link JournalEntry} that records actions by the {@link TaskManager}.
 *
 * <p>
 * Created Sep 17, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManagerJournalEntry extends BaseJournalEntry {

  private static String makeTitle(String taskName, TaskManagerActionType action) {
    return action.name() + ":" + taskName;
  }

  private static String makeDescription(String taskName, TaskManagerActionType action) {
    return makeTitle(taskName, action);
  }

  private final String taskName;
  private final TaskManagerActionType actionType;

  /**
   * Create a journal entry related to an action on a task.
   *
   * @param taskName The name of the task on which the action is taken
   * @param actionType The type of action that was taken on the task
   */
  public TaskManagerJournalEntry(String taskName, TaskManagerActionType actionType) {
    super(makeTitle(taskName, actionType), makeDescription(taskName, actionType));
    this.taskName = taskName;
    this.actionType = actionType;
  }

  public String getTaskName() {
    return taskName;
  }

  public TaskManagerActionType getActionType() {
    return actionType;
  }
}
