$(document).ready(function(){
    // 对于发布房源，只有认证后的用户才可以，所以先判断用户的实名认证状态
    $.get("/api/v1.0/user/auth", function(resp){
        if ("4101" == resp.errno) {
            // 用户未登录
            location.href = "/home/login.html";
        } else if ("0" == resp.errno) {
            // 未认证的用户，在页面中展示 "去认证"的按钮
            if (!(resp.data.real_name && resp.data.id_card)) {
                $(".auth-warn").show();
                return;
            }
            // 已认证的用户，请求其之前发布的房源信息
            $.get("/api/v1.0/user/houses", function(resp){
                if ("0" == resp.errno) {
                    $("#houses-list").html(template("houses-list-tmpl", {houses:resp.data.houses}));
                } else {
                    $("#houses-list").html(template("houses-list-tmpl", {houses:[]}));
                }
            });
        }
    });
    
})

function updateHouse(url,args) {
    console.log(args);
    location.href=url;
    event.stopPropagation( )
}
function del(e,id) {
    console.log(1);
    e.stopPropagation(); 
    console.log(id);
    if (confirm("确认要删除吗？")) {
        $.ajax({
        url: '/api/v1.0/houses/' + id,
        method: 'DELETE',
        success: function(resp) {
            if ("4101" == resp.errno) {
                location.href = "/home/login.html";
            } else if ("0" == resp.errno) {
                // 删除成功后的处理
                alert('删除成功！');
                location.reload(); // 刷新页面
            } else {
                alert(resp.errmsg);
            }
        },
        error: function(xhr, status, error) {
            // 删除失败后的处理
            alert('删除失败：' + error);
        }
        });
    }
}