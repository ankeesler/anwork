package com.marshmallow.anwork.task;

import java.util.Date;
import java.util.LinkedHashMap;
import java.util.Map;

import com.marshmallow.anwork.core.Serializer;

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


  private static class TaskSerializer implements Serializer<Task> {

    public static final Serializer<Task> instance = new TaskSerializer();

    private static final String START = "Task:";
    private static final String NAME = "name=";
    private static final String ID = "id=";
    private static final String DESCRIPTION = "description=";
    private static final String DATE = "date=";
    private static final String PRIORITY = "priority=";
    private static final String STATE = "state=";
    private static final String END = ";";

    @Override
    public String marshall(Task t) {
      StringBuilder builder = new StringBuilder();
      builder.append(START);
      builder.append(NAME).append(t.name).append(END);
      builder.append(ID).append(t.id).append(END);
      builder.append(DESCRIPTION).append(t.description).append(END);
      builder.append(DATE).append(t.startDate.toInstant().toEpochMilli()).append(END);
      builder.append(PRIORITY).append(t.priority).append(END);
      builder.append(STATE).append(t.state.name()).append(END);
      return builder.toString();
    }

    @Override
    public Task unmarshall(String string) {
      StringBuffer buffer = new StringBuffer(string);
      int index = buffer.indexOf(START);
      if (index != 0) {
        return null;
      }
      buffer.delete(index, START.length());

      Map<String, String> stuff = new LinkedHashMap<String, String>();
      stuff.put(NAME, "");
      stuff.put(ID, "");
      stuff.put(DESCRIPTION, "");
      stuff.put(DATE, "");
      stuff.put(PRIORITY, "");
      stuff.put(STATE, "");
      for (String key : stuff.keySet()) {
        int startIndex = buffer.indexOf(key);
        if (startIndex != 0) {
          return null;
        }
        int endIndex = buffer.indexOf(END, startIndex);
        if (endIndex == -1) {
          return null;
        }
        String value = buffer.substring(startIndex + key.length(), endIndex);
        stuff.put(key, value);
        buffer.delete(startIndex, endIndex + 1);
      }

      if (buffer.length() != 0) {
        return null;
      }

      Task task = new Task();
      try {
        task.name = stuff.get(NAME);
        task.id = Integer.parseInt(stuff.get(ID));
        task.description = stuff.get(DESCRIPTION);
        task.startDate = new Date(Long.parseLong(stuff.get(DATE)));
        task.priority = Integer.parseInt(stuff.get(PRIORITY));
        task.state = TaskState.valueOf(stuff.get(STATE));
      } catch (NumberFormatException nfe) {
        return null;
      }
      return task;
    }
  }

  /**
   * Get the instance {@link Serializer<Task>}.
   *
   * @return The instance {@link Serializer<Task>}
   */
  public static Serializer<Task> serializer() {
    return TaskSerializer.instance;
  }

  // This guy is here so that we can use him in the serialization functionality
  // above.
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
    return serializer().marshall(this);
  }
}
