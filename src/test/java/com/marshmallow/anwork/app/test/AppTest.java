package com.marshmallow.anwork.app.test;

import static org.junit.Assert.assertEquals;

import com.marshmallow.anwork.app.AnworkApp;
import com.marshmallow.anwork.core.FilePersister;
import com.marshmallow.anwork.core.Persister;
import com.marshmallow.anwork.core.test.TestUtilities;
import com.marshmallow.anwork.task.TaskManager;

import java.io.File;
import java.io.IOException;
import java.util.Collection;

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

  private void run(String...args) {
    AnworkApp.main(args);
  }

  @Test
  public void emptyArgsTest() {
    run();
  }

  @Test
  public void createTest() throws IOException {
    // FIXME: this is hardcoded based on internal FilePersister logic! Bad!
    File persistenceFile = new File(PERSISTENCE_ROOT, CONTEXT);
    if (persistenceFile.exists()) {
      persistenceFile.delete();
    }

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

  private TaskManager readTaskManager() throws IOException {
    Persister<TaskManager> persister = new FilePersister<TaskManager>(PERSISTENCE_ROOT);
    Collection<TaskManager> loadeds = persister.load(CONTEXT, TaskManager.serializer());
    assertEquals(1, loadeds.size());
    return loadeds.toArray(new TaskManager[0])[0];
  }
}
