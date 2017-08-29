package com.marshmallow.anwork.task;

import java.util.Date;

/**
 * This is a single unit of project work. One examples might be a single JIRA
 * ticket.
 *
 * @author Andrew
 * @date Aug 29, 2017
 */
public class Task implements Comparable<Task> {

  /** I totally made up this value. */
  public static int DEFAULT_PRIORITY = 5;

  private static int nextId = 1;

  private String name;
  private int id;
  private String description;
  private Date startDate;
  private int priority;
  private TaskState state;

  // There is only one constructor because we want to restrict the creation of
  // these objects to this package.
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
    builder.append("name=").append(name).append(';');
    builder.append("id=").append(id).append(';');
    builder.append("description=").append(description).append(';');
    builder.append("startDate=").append(startDate).append(';');
    builder.append("priority=").append(priority).append(';');
    builder.append("state=").append(state).append(';');
    return builder.toString();
  }
}
