package com.marshmallow.anwork.app.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;

import com.marshmallow.anwork.app.AnworkApp;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.core.test.TestUtilities;
import com.marshmallow.anwork.journal.Journal;
import com.marshmallow.anwork.task.Task;
import com.marshmallow.anwork.task.TaskManager;
import com.marshmallow.anwork.task.TaskState;

import java.io.File;
import java.io.IOException;
import java.util.Collection;

import org.junit.Before;
import org.junit.Test;

/**
 * This is a test for the anwork application.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class AppTest {

  private static final String CONTEXT = "app-test-context";
  private static final File PERSISTENCE_ROOT = TestUtilities.getFile(".", AppTest.class);

  /**
   * Hack-ish-ly clean up any existing context so we run our test cases fresh.
   */
  @Before
  public void removePreviousContext() throws IOException {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(PERSISTENCE_ROOT);
    if (persister.exists(CONTEXT)) {
      persister.clear(CONTEXT);
    }
  }

  @Test
  public void emptyArgsTest() {
    run();
  }

  @Test
  public void createTest() throws IOException {
    run("task", "create", "task-a",
        "-e", "This is the description for task A",
        "-p", "15");
    run("task", "create", "task-b",
        "--description", "This is the description for task B",
        "--priority", "25");
    run("task", "create", "task-c");

    TaskManager taskManager = readTaskManager();
    assertEquals(3, taskManager.getTasks().length);
    Task[] tasks = taskManager.getTasks();
    assertNotEquals(tasks[0].getId(), tasks[1].getId());
    assertNotEquals(tasks[0].getId(), tasks[2].getId());
    assertNotEquals(tasks[1].getId(), tasks[2].getId());
    assertEquals(3, taskManager.getJournal().getEntries().length);
    assertEquals(1, taskManager.getJournal("task-a").getEntries().length);
    assertEquals(1, taskManager.getJournal("task-b").getEntries().length);
    assertEquals(1, taskManager.getJournal("task-c").getEntries().length);

    Task taskA = null;
    Task taskB = null;
    for (Task task : taskManager.getTasks()) {
      if (task.getName().equals("task-a")) {
        taskA = task;
      } else if (task.getName().equals("task-b")) {
        taskB = task;
      }
    }
    assertNotNull("Unable to find task-a!", taskA);
    assertNotNull("Unable to find task-b!", taskB);
    assertEquals(15, taskA.getPriority());
    assertEquals("This is the description for task A", taskA.getDescription());
    assertEquals(25, taskB.getPriority());
    assertEquals("This is the description for task B", taskB.getDescription());
  }

  @Test
  public void setStateTest() throws IOException {
    run("task", "create", "task-a");
    run("task", "set-running", "task-a");
    run("task", "create", "task-b");
    run("task", "set-blocked", "task-b");

    TaskManager manager = readTaskManager();
    assertEquals(2, manager.getTasks().length);
    assertEquals(TaskState.RUNNING, manager.getState("task-a"));
    assertEquals(TaskState.BLOCKED, manager.getState("task-b"));
    assertEquals(4, manager.getJournal().getEntries().length);
    assertEquals(2, manager.getJournal("task-a").getEntries().length);
    assertEquals(2, manager.getJournal("task-b").getEntries().length);

    run("task", "set-finished", "task-a");

    manager = readTaskManager();
    assertEquals(2, manager.getTasks().length);
    assertEquals(TaskState.FINISHED, manager.getState("task-a"));
    assertEquals(TaskState.BLOCKED, manager.getState("task-b"));
    assertEquals(5, manager.getJournal().getEntries().length);
    assertEquals(3, manager.getJournal("task-a").getEntries().length);
    assertEquals(2, manager.getJournal("task-b").getEntries().length);

    run("task", "set-running", "task-b");

    manager = readTaskManager();
    assertEquals(2, manager.getTasks().length);
    assertEquals(TaskState.FINISHED, manager.getState("task-a"));
    assertEquals(TaskState.RUNNING, manager.getState("task-b"));
    assertEquals(6, manager.getJournal().getEntries().length);
    assertEquals(3, manager.getJournal("task-a").getEntries().length);
    assertEquals(3, manager.getJournal("task-b").getEntries().length);
  }

  @Test
  public void deleteTest() throws IOException {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "delete", "task-a");

    TaskManager taskManager = readTaskManager();
    assertEquals(1, taskManager.getTasks().length);
    assertEquals(3, taskManager.getJournal().getEntries().length);
    assertNull(taskManager.getJournal("task-a"));
    assertEquals(1, taskManager.getJournal("task-b").getEntries().length);

    run("task", "delete", "task-b");

    taskManager = readTaskManager();
    assertEquals(0, taskManager.getTasks().length);
    assertEquals(4, taskManager.getJournal().getEntries().length);
    assertNull(taskManager.getJournal("task-a"));
    assertNull(taskManager.getJournal("task-b"));
  }

  @Test
  public void deleteAllTest() throws IOException {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "create", "task-c");
    run("task", "delete-all");

    TaskManager taskManager = readTaskManager();
    assertEquals(0, taskManager.getTasks().length);
    assertEquals(6, taskManager.getJournal().getEntries().length);
  }

  @Test
  public void showTest() throws IOException {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "show");
    run("task", "delete", "task-a");
    run("task", "show");
  }

  @Test
  public void noteTest() throws IOException {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "create", "task-c");
    run("task", "note", "task-a", "hey task-a");
    run("task", "note", "task-c", "hey task-c");
    run("task", "note", "task-a", "hey task-a part 2");

    TaskManager taskManager = readTaskManager();
    assertEquals(3, taskManager.getTasks().length);
    assertEquals(6, taskManager.getJournal().getEntries().length);

    Journal<?> taskAJournal = taskManager.getJournal("task-a");
    assertNotNull("No journal for task-a!", taskAJournal);
    assertEquals(3, taskAJournal.getEntries().length);

    Journal<?> taskBJournal = taskManager.getJournal("task-b");
    assertNotNull("No journal for task-b!", taskBJournal);
    assertEquals(1, taskBJournal.getEntries().length);

    Journal<?> taskCJournal = taskManager.getJournal("task-c");
    assertNotNull("No journal for task-c!", taskCJournal);
    assertEquals(2, taskCJournal.getEntries().length);
  }

  @Test
  public void showAllJournalTest() throws IOException {
    run("journal", "show-all");
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("journal", "show-all");
    run("task", "delete", "task-a");
    run("journal", "show-all");
  }

  @Test
  public void showJournalTest() throws IOException {
    run("journal", "show", "task-a");
    run("task", "create", "task-a");
    run("journal", "show", "task-a");
    run("journal", "show", "task-b");
    run("task", "create", "task-b");
    run("journal", "show", "task-a");
    run("journal", "show", "task-b");
    run("task", "delete", "task-a");
    run("journal", "show", "task-a");
    run("journal", "show", "task-b");
    run("task", "set-finished", "task-b");
    run("journal", "show", "task-b");
  }

  @Test
  public void testNoPersist() throws IOException {
    run("--no-persist", "task", "create", "task-a");
    assertFalse(new FilePersister<TaskManager>(PERSISTENCE_ROOT).exists(CONTEXT));
  }

  private void run(String...args) {
    String[] baseArgs = new String[] {
      "-d",
      "--context", CONTEXT,
      "-o", PERSISTENCE_ROOT.getAbsolutePath(),
    };
    String[] allArgs = new String[baseArgs.length + args.length];
    System.arraycopy(baseArgs, 0, allArgs, 0, baseArgs.length);
    System.arraycopy(args, 0, allArgs, baseArgs.length, args.length);
    AnworkApp.main(allArgs);
  }

  private TaskManager readTaskManager() throws IOException {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(PERSISTENCE_ROOT);
    Collection<TaskManager> loadeds = persister.load(CONTEXT, TaskManager.SERIALIZER);
    assertEquals(1, loadeds.size());
    return loadeds.toArray(new TaskManager[0])[0];
  }
}
