package com.marshmallow.anwork.task;

import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.journal.MultiJournaled;
import com.marshmallow.anwork.task.protobuf.TaskManagerProtobuf;
import com.marshmallow.anwork.task.protobuf.TaskProtobuf;

import java.io.IOException;
import java.util.Arrays;
import java.util.PriorityQueue;
import java.util.function.Consumer;
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
    journal.addEntry(new TaskManagerJournalEntry(task, TaskManagerActionType.CREATE, ""));
  }

  /**
   * Delete a task from a name.
   *
   * @param name The name for the task
   * @throws IllegalArgumentException If this task does not exist.
   */
  public void deleteTask(String name) throws IllegalArgumentException {
    doWithTask(name, (task) -> {
      tasks.remove(task);
      journal.addEntry(new TaskManagerJournalEntry(task, TaskManagerActionType.DELETE, ""));
    });
  }

  /**
   * Get the state of a task.
   *
   * @param name The name of the task
   * @return A string representing the state of the task.
   * @throws IllegalArgumentException If the task does not exist
   */
  public TaskState getState(String name) throws IllegalArgumentException {
    final TaskState[] stateHolder = new TaskState[1];
    doWithTask(name, (task) -> {
      stateHolder[0] = task.getState();
    });
    return stateHolder[0];
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
    doWithTask(name, (task) -> {
      task.setState(state);
      journal.addEntry(new TaskManagerJournalEntry(task,
                                                   TaskManagerActionType.SET_STATE,
                                                   state.name()));
    });
  }

  /**
   * Add a note to a task.
   *
   * @param name The name of the task
   * @param note The note to add to the task
   * @throws IllegalArgumentException If this task does not exist
   */
  public void addNote(String name, String note) throws IllegalArgumentException {
    doWithTask(name, (task) -> {
      journal.addEntry(new TaskManagerJournalEntry(task,
                                                   TaskManagerActionType.NOTE,
                                                   note));
    });
  }

  /**
   * Set a {@link Task}'s priority.
   *
   * @param name The name of the task
   * @param priority The priority to set on the task
   * @throws IllegalArgumentException If this task does not exist
   */
  public void setPriority(String name, int priority) throws IllegalArgumentException {
    doWithTask(name, (task) -> {
      journal.addEntry(new TaskManagerJournalEntry(task,
                                                   TaskManagerActionType.SET_PRIORITY,
                                                   Integer.toString(priority)));
      task.setPriority(priority);
    });
  }

  /**
   * Get the {@link Task}'s associated with this {@link TaskManager}. The {@link Task}'s are
   * returned sorted by priority (lowest priority number to highest priority number).
   *
   * @return The {@link Task}'s associated with this {@link TaskManager}
   */
  public Task[] getTasks() {
    // Eh, this is a bummer that we have to sort an array that should really already be sorted.
    // But, Java's PriorityQueue doesn't do this for us by default. We could write our own, but
    // that seems silly.
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

  private void doWithTask(String name, Consumer<Task> doer) {
    Task task = findTask(name);
    if (task == null) {
      throw new IllegalArgumentException("Task " + name + " does not exist");
    }
    doer.accept(task);
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
  public TaskManagerProtobuf marshall() throws IOException {
    TaskManagerProtobuf.Builder builder = TaskManagerProtobuf.newBuilder();
    Task[] taskArray = tasks.toArray(new Task[0]);
    for (Task task : taskArray) {
      builder.addTasks(task.marshall());
    }
    builder.setJournal(journal.marshall());
    return builder.build();
  }

  @Override
  public void unmarshall(TaskManagerProtobuf protobuf) throws IOException {
    for (int i = 0; i < protobuf.getTasksCount(); i++) {
      TaskProtobuf taskProtobuf = protobuf.getTasks(i);
      Task task = Task.FACTORY.makeBlankInstance();
      task.unmarshall(taskProtobuf);
      tasks.add(task);
    }
    journal.unmarshall(protobuf.getJournal());
  }

  @Override
  public Journal<TaskManagerJournalEntry> getJournal() {
    return journal;
  }

  /**
   * Get the {@link Journal} associated with a {@link Task}.
   *
   * @param key The name of the {@link Task} for which to get the {@link Journal}
   * @return The journal associated with a {@link Task}, or <code>null</code> if there is no
   *     {@link Task} with the provided name
   */
  @Override
  public Journal<TaskManagerJournalEntry> getJournal(String key) {
    return journal.filter(key);
  }
}