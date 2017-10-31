package com.marshmallow.anwork.task;

import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.journal.JournalEntry;
import com.marshmallow.anwork.task.protobuf.TaskManagerActionTypeProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskManagerJournalEntryProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskProtobuf;

import java.io.IOException;
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
public class TaskManagerJournalEntry implements JournalEntry,
                                                Serializable<TaskManagerJournalEntryProtobuf> {

  public static final Serializer<TaskManagerJournalEntry> SERIALIZER
      = new ProtobufSerializer<TaskManagerJournalEntryProtobuf,
                               TaskManagerJournalEntry>(() -> new TaskManagerJournalEntry(),
                                                        TaskManagerJournalEntryProtobuf.parser());

  private Task task;
  private TaskManagerActionType actionType;
  private String detail;
  private Date date;

  // Default constructor for the factory in SERIALIZER and use in TaskManagerJournal.
  TaskManagerJournalEntry() {
  }

  /**
   * Create a journal entry related to a {@link TaskManagerActionType} of action on a {@link Task}.
   *
   * @param task The {@link Task} on which the action of type {@link TaskManagerActionType} is
   *     taken
   * @param actionType The {@link TaskManagerActionType} of action that was taken on the
   *     {@link Task}
   * @param detail A {@link String} containing detail of the event related to a {@link Task}
   */
  public TaskManagerJournalEntry(Task task, TaskManagerActionType actionType, String detail) {
    this.task = task;
    this.actionType = actionType;
    this.detail = detail;
    this.date = new Date();
  }

  @Override
  public String getTitle() {
    switch (actionType) {
      case CREATE:
        return String.format("Created task '%s'", task.getName());
      case DELETE:
        return String.format("Deleted task '%s'", task.getName());
      case SET_STATE:
        return String.format("Set %s state on task '%s'", detail, task.getName());
      case NOTE:
        return String.format("Added note to task '%s': %s", task.getName(), detail);
      default:
        return null;
    }
  }

  @Override
  public String getDescription() {
    return String.format("Action %s was performed on task '%s'. Detail: %s",
                         actionType.name(),
                         task.getName(),
                         detail);
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

  /**
   * Get the detail {@link String} for this {@link TaskManagerJournalEntry}.
   *
   * @return The detail {@link String} for this {@link TaskManagerJournalEntry}
   */
  public String getDetail() {
    return detail;
  }

  @Override
  public TaskManagerJournalEntryProtobuf marshall() throws IOException {
    TaskProtobuf taskProtobuf = task.marshall();
    TaskManagerJournalEntryProtobuf.Builder builder = TaskManagerJournalEntryProtobuf.newBuilder();
    builder.setTask(taskProtobuf);
    builder.setActionType(TaskManagerActionTypeProtobuf.forNumber(actionType.ordinal()));
    builder.setDetail(detail);
    builder.setDate(date.getTime());
    return builder.build();
  }

  @Override
  public void unmarshall(TaskManagerJournalEntryProtobuf t) throws IOException {
    task = Task.FACTORY.makeBlankInstance();
    task.unmarshall(t.getTask());
    actionType = TaskManagerActionType.values()[t.getActionType().ordinal()];
    detail = t.getDetail();
    date = new Date(t.getDate());
  }
}
