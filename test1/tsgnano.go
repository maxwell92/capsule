package main
 
import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    "time"
 
    "github.com/docker/docker/pkg/mount"
    "github.com/docker/docker/pkg/reexec"
)
 
func init() {
    reexec.Register("namespaced", namespaced)
 
    if reexec.Init() {
        os.Exit(0)
    }
}
 
func namespaced() {
    id := os.Args[1]
    fmt.Printf("RUNNING IN NEW NAMESPACE WITH ID: %s", id)
    pivot_root_dir := fmt.Sprintf("/tmp/pivot-%s", id)
 
    // pivot root here
    if err := os.MkdirAll(pivot_root_dir, 0755); err != nil {
        fmt.Printf("failed to mkdir: %s\n", err)
        os.Exit(1)
    }
 
    if err := pivotRoot(pivot_root_dir); err != nil {
        fmt.Printf("failed to pivot_root: %s\n", err)
        os.Exit(1)
    }
 
    time.Sleep(time.Hour)
 
    // programPath, err := exec.LookPath(os.Args[1])
    // if err != nil {
    //  fmt.Printf("failed to lookup path in namespace : %s\n", err)
    //  os.Exit(1)
    // }
 
    // err = syscall.Exec(programPath, os.Args[1:], os.Environ())
    // if err != nil {
    //  fmt.Printf("exec failed in namespace: %s\n", err)
    //  os.Exit(1)
    // }
}
 
func main() {
    reexecArgs := append([]string{"namespaced"}, os.Args[1:]...)
    cmd := reexec.Command(reexecArgs...)
 
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
    }
 
    if err := cmd.Run(); err != nil {
        fmt.Printf("tsgnano failed: %s\n", err)
        os.Exit(1)
    }
}
 
func pivotRoot(rootfs string) error {
    // While the documentation may claim otherwise, pivot_root(".", ".") is
    // actually valid. What this results in is / being the new root but
    // /proc/self/cwd being the old root. Since we can play around with the cwd
    // with pivot_root this allows us to pivot without creating directories in
    // the rootfs. Shout-outs to the LXC developers for giving us this idea.
 
    oldroot, err := syscall.Open("/", syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
    if err != nil {
        return err
    }
    defer syscall.Close(oldroot)
 
    newroot, err := syscall.Open(rootfs, syscall.O_DIRECTORY|syscall.O_RDONLY, 0)
    if err != nil {
        return err
    }
    defer syscall.Close(newroot)
 
    // Change to the new root so that the pivot_root actually acts on it.
    if err := syscall.Fchdir(newroot); err != nil {
        return err
    }
 
    if err := syscall.PivotRoot(".", "."); err != nil {
        // Make the parent mount private
        if err := rootfsParentMountPrivate("."); err != nil {
            return err
        }
 
        // Try again
        if err := syscall.PivotRoot(".", "."); err != nil {
            return fmt.Errorf("pivot_root %s", err)
        }
    }
 
    // Currently our "." is oldroot (according to the current kernel code).
    // However, purely for safety, we will fchdir(oldroot) since there isn't
    // really any guarantee from the kernel what /proc/self/cwd will be after a
    // pivot_root(2).
 
    if err := syscall.Fchdir(oldroot); err != nil {
        return err
    }
 
    // Make oldroot rprivate to make sure our unmounts don't propogate to the
    // host (and thus bork the machine).
    if err := syscall.Mount("", ".", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
        return err
    }
    // Preform the unmount. MNT_DETACH allows us to unmount /proc/self/cwd.
    if err := syscall.Unmount(".", syscall.MNT_DETACH); err != nil {
        return err
    }
 
    // Switch back to our shiny new root.
    if err := syscall.Chdir("/"); err != nil {
        return fmt.Errorf("chdir / %s", err)
    }
    return nil
}
 
func rootfsParentMountPrivate(rootfs string) error {
    sharedMount := false
 
    parentMount, optionalOpts, err := getParentMount(rootfs)
    if err != nil {
        return err
    }
 
    optsSplit := strings.Split(optionalOpts, " ")
    for _, opt := range optsSplit {
        if strings.HasPrefix(opt, "shared:") {
            sharedMount = true
            break
        }
    }
 
    // Make parent mount PRIVATE if it was shared. It is needed for two
    // reasons. First of all pivot_root() will fail if parent mount is
    // shared. Secondly when we bind mount rootfs it will propagate to
    // parent namespace and we don't want that to happen.
    if sharedMount {
        return syscall.Mount("", parentMount, "", syscall.MS_PRIVATE, "")
    }
 
    return nil
}
 
func getParentMount(rootfs string) (string, string, error) {
    var path string
 
    mountinfos, err := mount.GetMounts()
    if err != nil {
        return "", "", err
    }
 
    mountinfo := getMountInfo(mountinfos, rootfs)
    if mountinfo != nil {
        return rootfs, mountinfo.Optional, nil
    }
 
    path = rootfs
    for {
        path = filepath.Dir(path)
 
        mountinfo = getMountInfo(mountinfos, path)
        if mountinfo != nil {
            return path, mountinfo.Optional, nil
        }
 
        if path == "/" {
            break
        }
    }
 
    // If we are here, we did not find parent mount. Something is wrong.
    return "", "", fmt.Errorf("Could not find parent mount of %s", rootfs)
}
 
func getMountInfo(mountinfo []*mount.Info, dir string) *mount.Info {
    for _, m := range mountinfo {
        if m.Mountpoint == dir {
            return m
        }
    }
    return nil
}
