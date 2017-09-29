package com.marshmallow.anwork.task;

import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.journal.MultiJournaled;
import com.marshmallow.anwork.task.protobuf.TaskManagerProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskProtobuf;

import java.util.Arrays;
import java.util.PriorityQueue;
import java.util.stream.Stream;

/**
 * This guy in a public interface for managing {@link Task} instances.
 *
 * <p>
 * Created Aug 29, 2017
 * </p>
 *
 * @author Andrew
 */
public class TaskManager implements Serializable<TaskManagerProtobuf>,
                                    MultiJournaled<TaskManagerJournalEntry> {

  /**
   * This is the singleton {@link Serializer} for this class.
   */
  public static Serializer<TaskManager> SERIALIZER
      = new ProtobufSerializer<TaskManagerProtobuf, TaskManager>(() -> new TaskManager(),
                                                                 TaskManagerProtobuf.parser());

  private PriorityQueue<Task> tasks = new PriorityQueue<Task>();

  private TaskManagerJournal journal = new TaskManagerJournal();

  /**
   * Create a task from a name.
   *
   * @param name The name for the task
   * @param description The description for the task
   * @param priority The priority for the task
   * @throws IllegalArgumentException If this task already exists
   */
  public void createTask(String name,
                         String description,
                         int priority) throws IllegalArgumentException {
    if (findTask(name) != null) {
      throw new IllegalArgumentException("Task " + name + " already exists");
    }
    Task task = new Task(name, description, priority);
    tasks.add(task);
    journal.addEntry(new TaskManagerJournalEntry(task, TaskManagerActionType.CREATE));
  }

  /**
   * Delete a task from a name.
   *
   * @param name The name for the task
   * @throws IllegalArgumentException If this task does not exist.
   */
  public void deleteTask(String name) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    tasks.remove(task);
    journal.addEntry(new TaskManagerJournalEntry(task, TaskManagerActionType.DELETE));
  }

  /**
   * Get the state of a task.
   *
   * @param name The name of the task
   * @return A string representing the state of the task.
   * @throws IllegalArgumentException If the task does not exist
   */
  public TaskState getState(String name) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    return task.getState();
  }

  /**
   * Set the state of a task.
   *
   * @param name The name of the task
   * @param state The name of the state.
   * @throws IllegalArgumentException If this task does not exist or the state
   *     is invalid.
   */
  public void setState(String name, TaskState state) throws IllegalArgumentException {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    task.setState(state);
    journal.addEntry(new TaskManagerJournalEntry(task, TaskManagerActionType.SET_STATE));
  }

  /**
   * Get the tasks associated with this task manager.
   *
   * @return The tasks associated with this manager
   */
  public Task[] getTasks() {
    // TODO: this should be improved! We should write our own heap that can
    // return the elements to us in order.
    Task[] taskArray = tasks.toArray(new Task[0]);
    Arrays.sort(taskArray);
    return taskArray;
  }

  /**
   * Get a string representation of the tasks in this manager.
   * @return A string representation of the tasks in this manager.
   */
  @Override
  public String toString() {
    StringBuilder builder = new StringBuilder();
    Stream.of(tasks.toArray(new Task[0])).forEach(task -> builder.append(task));
    return builder.toString();
  }

  private Task findTask(String name) {
    Task[] taskArray = tasks.toArray(new Task[0]);
    for (Task task : taskArray) {
      if (task.getName().equals(name)) {
        return task;
      }
    }
    return null;
  }

  @Override
  public TaskManagerProtobuf marshall() {
    TaskManagerProtobuf.Builder builder = TaskManagerProtobuf.newBuilder();
    Task[] taskArray = tasks.toArray(new Task[0]);
    for (Task task : taskArray) {
      builder.addTasks(task.marshall());
    }
    return builder.build();
  }

  @Override
  public void unmarshall(TaskManagerProtobuf protobuf) {
    for (int i = 0; i < protobuf.getTasksCount(); i++) {
      TaskProtobuf taskProtobuf = protobuf.getTasks(i);
      Task task = Task.FACTORY.makeBlankInstance();
      task.unmarshall(taskProtobuf);
      tasks.add(task);
    }
  }

  @Override
  public Journal<TaskManagerJournalEntry> getJournal() {
    return journal;
  }

  /**
   * Get the {@link Journal} associated with a {@link Task}.
   *
   * @param key The name of the {@link Task} for which to get the {@link Journal}
   * @return The journal associated with a {@link Task}, or null if there is no {@link Task} with
   *     the provided name
   */
  @Override
  public Journal<TaskManagerJournalEntry> getJournal(String key) {
    return journal.filter(key);
  }
}