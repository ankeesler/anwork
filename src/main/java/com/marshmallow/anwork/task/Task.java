package com.marshmallow.anwork.task;

import com.marshmallow.anwork.core.Factory;
import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.journal.BaseJournal;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.journal.Journaled;
import com.marshmallow.anwork.task.protobuf.TaskProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskStateProtobuf;

import java.util.Date;

/**
 * This is a single unit of project work. One examples might be a single JIRA
 * ticket.
 *
 * <p>
 * Created Aug 29, 2017
 * </p>
 *
 * @author Andrew
 */
public class Task implements Comparable<Task>, Serializable<TaskProtobuf>, Journaled {

  /** I totally made up this value. */
  public static int DEFAULT_PRIORITY = 5;

  /**
   * This is the singleton {@link Factory} for this class. This is meant to only be used
   * in this package.
   */
  static Factory<Task> FACTORY = () -> new Task();

  /**
   * This is the singleton {@link Serializer} for this class.
   */
  public static Serializer<Task> SERIALIZER
      = new ProtobufSerializer<TaskProtobuf, Task>(FACTORY, TaskProtobuf.parser());

  private static int nextId = 1;

  private String name;
  private int id;
  private String description;
  private Date startDate;
  private int priority;
  private TaskState state;

  private BaseJournal journal;

  // This is for the factory above.
  private Task() { }

  // There is only one non-private constructor because we want to restrict the
  // creation of these objects to this package.
  Task(String name, String description, int priority) {
    this.name = name;
    this.id = nextId++;
    this.description = description;
    this.startDate = new Date();
    this.priority = priority;
    this.state = TaskState.WAITING;
  }

  public TaskState getState() {
    return state;
  }

  public void setState(TaskState state) {
    this.state = state;
  }

  public String getName() {
    return name;
  }

  public int getId() {
    return id;
  }

  public String getDescription() {
    return description;
  }

  public Date getStartDate() {
    return startDate;
  }

  public int getPriority() {
    return priority;
  }

  @Override
  public int compareTo(Task o) {
    return this.priority - o.priority;
  }

  @Override
  public int hashCode() {
    return name.hashCode();
  }

  @Override
  public boolean equals(Object o) {
    if (o == null) {
      return this == null;
    } else if (!(o instanceof Task)) {
      return false;
    } else if (((Task)o).name == null) {
      return name == null;
    } else {
      return name.equals(((Task)o).name);
    }
  }

  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder();
    builder.append("Task:");
    builder.append(" n='").append(name).append("'");
    builder.append(" d='").append(description).append("'");
    builder.append(" s='").append(startDate.toString()).append("'");
    builder.append(" p='").append(priority).append("'");
    builder.append(" t='").append(state).append("'");
    return builder.toString();
  }

  @Override
  public TaskProtobuf marshall() {
    return TaskProtobuf.newBuilder()
                       .setId(id)
                       .setName(name)
                       .setDescription(description)
                       .setStartDate(startDate.toInstant().getEpochSecond())
                       .setPriority(priority)
                       .setState(TaskStateProtobuf.forNumber(state.ordinal()))
                       .build();
  }

  @Override
  public void unmarshall(TaskProtobuf t) {
    id = t.getId();
    name = t.getName();
    description = t.getDescription();
    startDate = new Date(t.getStartDate());
    priority = t.getPriority();
    state = TaskState.values()[t.getState().ordinal()];
  }

  @Override
  public Journal getJournal() {
    return journal;
  }
}
