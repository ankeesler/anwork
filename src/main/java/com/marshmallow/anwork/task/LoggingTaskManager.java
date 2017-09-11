package com.marshmallow.anwork.task;

import com.marshmallow.anwork.event.Event;
import com.marshmallow.anwork.event.EventLog;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * This is a {@link TaskManager} that logs its CRUD operations.
 *
 * @author Andrew
 * Created Aug 31, 2017
 */
public class LoggingTaskManager extends TaskManager {

  private List<EventLog> logs = new ArrayList<EventLog>();

  private static class TaskManagerEvent implements Event {

    private String type;
    private String taskName;
    private Date date;
    private boolean success;

    public TaskManagerEvent(String type,
                            String taskName,
                            boolean success) {
      this.type = type;
      this.taskName = taskName;
      this.date = new Date();
      this.success = success;
    }

    @Override
    public String getType() {
      return type;
    }

    @Override
    public Date getDate() {
      return date;
    }

    @Override
    public String getDescription() {
      return String.format("%s task %s (%s)",
                           type,
                           taskName,
                           (success ? "success" : "failure"));
    }
    
  }

  @Override
  public void createTask(String name,
                         String description,
                         int priority) throws IllegalArgumentException {
    boolean success = true;
    try {
      super.createTask(name, description, priority);
    } catch (IllegalArgumentException iae) {
      success = false;
      throw iae;
    } finally {
      Event event = new TaskManagerEvent("create", name, success);
      logs.forEach(log -> log.add(event));
    }
  }

  @Override
  public void deleteTask(String name) throws IllegalArgumentException {
    boolean success = true;
    try {
      super.deleteTask(name);
    } catch (IllegalArgumentException iae) {
      success = false;
      throw iae;
    } finally {
      Event event = new TaskManagerEvent("delete", name, success);
      logs.forEach(log -> log.add(event));
    }
  }

  @Override
  public void setCurrentTask(String name) throws IllegalArgumentException {
    boolean success = true;
    try {
      super.setCurrentTask(name);
    } catch (IllegalArgumentException iae) {
      success = false;
      throw iae;
    } finally {
      Event event = new TaskManagerEvent("set-current-task", name, success);
      logs.forEach(log -> log.add(event));
    }
  }

  @Override
  public void setState(String name, TaskState state) throws IllegalArgumentException {
    boolean success = true;
    try {
      super.setState(name, state);
    } catch (IllegalArgumentException iae) {
      success = false;
      throw iae;
    } finally {
      Event event = new TaskManagerEvent("set-state", name, success);
      logs.forEach(log -> log.add(event));
    }
  }

  /**
   * Add an {@link EventLog} to this manager.
   *
   * The log will be used to record the manager's CRUD operations.
   *
   * @param log The log to add
   */
  public void addLog(EventLog log) {
    logs.add(log);
  }

  /**
   * Remove a log from this manager.
   *
   * @param log The log to remove
   * @return Whether or not a log was removed
   */
  public boolean removeLog(EventLog log) {
    return logs.remove(log);
  }
}
