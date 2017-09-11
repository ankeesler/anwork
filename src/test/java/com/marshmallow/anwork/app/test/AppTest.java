package com.marshmallow.anwork.app.test;

import static org.junit.Assert.assertEquals;

import com.marshmallow.anwork.app.AnworkApp;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.core.test.TestUtilities;
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
 * @author Andrew
 * Created Sep 9, 2017
 */
public class AppTest {

  private static final String CONTEXT = "app-test-context";
  private static final File PERSISTENCE_ROOT
    = new File(TestUtilities.TEST_RESOURCES_ROOT, "app-test");

  @Before
  public void removePreviousContext() {
    // FIXME: this is hardcoded based on internal FilePersister logic! Bad!
    File persistenceFile = new File(PERSISTENCE_ROOT, CONTEXT);
    if (persistenceFile.exists()) {
      persistenceFile.delete();
    }
  }

  @Test
  public void emptyArgsTest() {
    run();
  }

  @Test
  public void createTest() throws IOException {
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "create", "task-a", "This is the description for task A", "1");
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "create", "task-b", "This is the description for task B", "2");
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "create", "task-c", "This is the description for task C", "3");

    TaskManager taskManager = readTaskManager();
    assertEquals(3, taskManager.getTaskCount());
  }

  @Test
  public void setStateTest() throws IOException {
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "create", "task-a", "This is the description for task A", "1");
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "set-running", "task-a");
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "create", "task-b", "This is the description for task B", "2");
    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "set-blocked", "task-b");

    TaskManager manager = readTaskManager();
    assertEquals(2, manager.getTaskCount());
    assertEquals(TaskState.RUNNING, manager.getState("task-a"));
    assertEquals(TaskState.BLOCKED, manager.getState("task-b"));

    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "set-finished", "task-a");

    manager = readTaskManager();
    assertEquals(2, manager.getTaskCount());
    assertEquals(TaskState.FINISHED, manager.getState("task-a"));
    assertEquals(TaskState.BLOCKED, manager.getState("task-b"));

    run("-d",
        "--context", CONTEXT,
        "-o", PERSISTENCE_ROOT.getAbsolutePath(),
        "task", "set-running", "task-b");

    manager = readTaskManager();
    assertEquals(2, manager.getTaskCount());
    assertEquals(TaskState.FINISHED, manager.getState("task-a"));
    assertEquals(TaskState.RUNNING, manager.getState("task-b"));
  }

  private void run(String...args) {
    AnworkApp.main(args);
  }

  private TaskManager readTaskManager() throws IOException {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(PERSISTENCE_ROOT);
    Collection<TaskManager> loadeds = persister.load(CONTEXT, TaskManager.serializer());
    assertEquals(1, loadeds.size());
    return loadeds.toArray(new TaskManager[0])[0];
  }
}
